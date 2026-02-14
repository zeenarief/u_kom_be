package service

import (
	"time"
	"u_kom_be/internal/apperrors"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
	"u_kom_be/internal/utils"
)

type AttendanceService interface {
	SubmitAttendance(req request.AttendanceSubmitRequest) (*response.AttendanceSessionDetailResponse, error)
	GetSessionDetail(id string) (*response.AttendanceSessionDetailResponse, error)
	GetHistoryByTeacher(teacherID string) ([]response.AttendanceHistoryResponse, error)
	GetHistoryByAssignment(taID string) ([]response.AttendanceHistoryResponse, error)
	GetSessionByScheduleDate(scheduleID, dateStr string) (*response.AttendanceSessionDetailResponse, error)
	GetSessionOrClassList(scheduleID, dateStr string) (*response.AttendanceSessionDetailResponse, error)
	DeleteSession(id string) error
}

type attendanceService struct {
	repo         repository.AttendanceRepository
	scheduleRepo repository.ScheduleRepository
	studentRepo  repository.StudentRepository
}

func NewAttendanceService(
	repo repository.AttendanceRepository,
	schedRepo repository.ScheduleRepository,
	studRepo repository.StudentRepository,
) AttendanceService {
	return &attendanceService{
		repo:         repo,
		scheduleRepo: schedRepo,
		studentRepo:  studRepo,
	}
}

func (s *attendanceService) SubmitAttendance(req request.AttendanceSubmitRequest) (*response.AttendanceSessionDetailResponse, error) {
	// PERBAIKAN: Parse menggunakan Local Location server
	date, err := time.ParseInLocation("2006-01-02", req.Date, time.Local)
	if err != nil {
		return nil, apperrors.NewBadRequestError("invalid date format")
	}

	// Siapkan detail baru
	var newDetails []domain.AttendanceDetail
	for _, studentInput := range req.Students {
		newDetails = append(newDetails, domain.AttendanceDetail{
			StudentID: studentInput.StudentID,
			Status:    studentInput.Status,
			Notes:     studentInput.Notes,
		})
	}

	// 1. CEK EKSISTENSI (PERBAIKAN UTAMA DISINI)
	// Jangan gunakan underscore (_), tangkap errornya
	existingSession, errFind := s.repo.FindSessionByScheduleDate(req.ScheduleID, date)

	// Logic: Jika tidak ada error DAN session ditemukan -> UPDATE
	if errFind == nil && existingSession != nil && existingSession.ID != "" {

		// === UPDATE MODE ===
		existingSession.Topic = req.Topic
		existingSession.Notes = req.Notes

		// Panggil repo Update
		if err := s.repo.UpdateSession(existingSession, newDetails); err != nil {
			return nil, err
		}

		return s.GetSessionDetail(existingSession.ID)

	} else {
		// === CREATE MODE ===
		// Masuk sini jika record not found (errFind != nil)
		session := &domain.AttendanceSession{
			ScheduleID: req.ScheduleID,
			Date:       utils.Date(date),
			Topic:      req.Topic,
			Notes:      req.Notes,
			Details:    newDetails,
		}

		if err := s.repo.CreateSession(session); err != nil {
			return nil, err
		}

		return s.GetSessionDetail(session.ID)
	}
}

// GetSessionDetail mengambil detail sesi + rekap kehadiran
func (s *attendanceService) GetSessionDetail(id string) (*response.AttendanceSessionDetailResponse, error) {
	session, err := s.repo.FindSessionByID(id)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, apperrors.NewNotFoundError("attendance session not found")
	}

	// Mapping ke Response
	res := &response.AttendanceSessionDetailResponse{
		ID:    session.ID,
		Date:  session.Date.Format("2006-01-02"),
		Topic: session.Topic,
		Notes: session.Notes,
		// Map info jadwal ringkas
		ScheduleInfo: response.ScheduleResponse{
			ID:            session.Schedule.ID,
			DayOfWeek:     session.Schedule.DayOfWeek,
			StartTime:     session.Schedule.StartTime,
			EndTime:       session.Schedule.EndTime,
			SubjectName:   session.Schedule.TeachingAssignment.Subject.Name,
			ClassroomName: session.Schedule.TeachingAssignment.Classroom.Name,
			TeacherName:   session.Schedule.TeachingAssignment.Teacher.FullName,
		},
		Summary: make(map[string]int),
	}

	// Loop details untuk isi list siswa & hitung summary
	for _, d := range session.Details {
		// Add to details list
		res.Details = append(res.Details, response.AttendanceDetailResponse{
			StudentID:   d.StudentID,
			StudentName: d.Student.FullName, // Asumsi Preload Student berhasil
			NISN:        utils.SafeString(d.Student.NISN),
			Status:      d.Status,
			Notes:       d.Notes,
		})

		// Increment Summary
		res.Summary[d.Status]++
	}

	return res, nil
}

// GetHistoryByTeacher menampilkan riwayat mengajar guru tertentu
func (s *attendanceService) GetHistoryByTeacher(teacherID string) ([]response.AttendanceHistoryResponse, error) {
	sessions, err := s.repo.GetHistoryByTeacher(teacherID)
	if err != nil {
		return nil, err
	}

	var history []response.AttendanceHistoryResponse
	for _, sess := range sessions {
		// Hitung jumlah yang TIDAK HADIR (Sakit/Izin/Alpa) - Opsional logic
		// Disini kita perlu preload details di repo GetHistoryByTeacher jika ingin hitung akurat
		// Jika query repo belum preload details, count_absent akan 0.
		// Untuk performa list, biasanya count dilakukan di query SQL, tapi untuk sekarang kita skip atau biarkan 0.

		history = append(history, response.AttendanceHistoryResponse{
			ID:          sess.ID,
			Date:        sess.Date,
			SubjectName: sess.Schedule.TeachingAssignment.Subject.Name,
			ClassName:   sess.Schedule.TeachingAssignment.Classroom.Name,
			Topic:       sess.Topic,
			CountAbsent: 0, // Placeholder, butuh query count spesifik jika ingin ditampilkan di list
		})
	}

	return history, nil
}

func (s *attendanceService) GetHistoryByAssignment(taID string) ([]response.AttendanceHistoryResponse, error) {
	sessions, err := s.repo.GetHistoryByTeachingAssignmentID(taID)
	if err != nil {
		return nil, err
	}

	var history []response.AttendanceHistoryResponse
	for _, sess := range sessions {
		// fmt.Printf("DEBUG ID: %s, ScheduleID: %s\n", sess.ID, sess.ScheduleID)
		history = append(history, response.AttendanceHistoryResponse{
			ID:          sess.ID,
			Date:        sess.Date,
			ScheduleID:  sess.ScheduleID,
			SubjectName: sess.Schedule.TeachingAssignment.Subject.Name,
			ClassName:   sess.Schedule.TeachingAssignment.Classroom.Name,
			Topic:       sess.Topic,
			CountAbsent: 0, // Placeholder
		})
	}
	return history, nil
}

func (s *attendanceService) GetSessionByScheduleDate(scheduleID, dateStr string) (*response.AttendanceSessionDetailResponse, error) {
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		return nil, err
	}

	session, err := s.repo.FindSessionByScheduleDate(scheduleID, date)
	if err != nil || session == nil {
		return nil, apperrors.NewNotFoundError("not found")
	}

	// Reuse logic mapping
	return s.GetSessionDetail(session.ID)
}

// GetSessionOrClassList mengembalikan sesi yang ada ATAU list siswa jika belum ada absen
func (s *attendanceService) GetSessionOrClassList(scheduleID, dateStr string) (*response.AttendanceSessionDetailResponse, error) {
	// 1. Cek apakah sesi absen sudah ada
	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		return nil, apperrors.NewBadRequestError("invalid date format")
	}

	session, err := s.repo.FindSessionByScheduleDate(scheduleID, date)
	if err == nil && session != nil {
		// Jika ADA, kembalikan detail sesi
		return s.GetSessionDetail(session.ID)
	}

	// 2. Jika BELUM ADA, ambil data Schedule -> Classroom -> Students
	schedule, err := s.scheduleRepo.FindByID(scheduleID)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, apperrors.NewNotFoundError("schedule not found")
	}

	// Ambil list siswa yang aktif di kelas tersebut
	students, err := s.studentRepo.FindByClassroomID(schedule.TeachingAssignment.ClassroomID)
	if err != nil {
		return nil, err
	}

	// 3. Construct response "kosong" tapi berisi list siswa
	res := &response.AttendanceSessionDetailResponse{
		ID:    "", // Kosong menandakan belum disave
		Date:  dateStr,
		Topic: "",
		Notes: "",
		ScheduleInfo: response.ScheduleResponse{
			ID:            schedule.ID,
			DayOfWeek:     schedule.DayOfWeek,
			StartTime:     schedule.StartTime,
			EndTime:       schedule.EndTime,
			SubjectName:   schedule.TeachingAssignment.Subject.Name,
			ClassroomName: schedule.TeachingAssignment.Classroom.Name,
			TeacherName:   schedule.TeachingAssignment.Teacher.FullName,
		},
		Details: []response.AttendanceDetailResponse{},
		Summary: make(map[string]int),
	}

	for _, student := range students {
		res.Details = append(res.Details, response.AttendanceDetailResponse{
			StudentID:   student.ID,
			StudentName: student.FullName,
			NISN:        utils.SafeString(student.NISN),
			Status:      "", // Kosong atau default "PRESENT"
			Notes:       "",
		})
	}

	return res, nil
}

func (s *attendanceService) DeleteSession(id string) error {
	// Cek apakah session ada
	session, err := s.repo.FindSessionByID(id)
	if err != nil {
		return err
	}
	if session == nil {
		return apperrors.NewNotFoundError("attendance session not found")
	}

	return s.repo.DeleteSession(id)
}
