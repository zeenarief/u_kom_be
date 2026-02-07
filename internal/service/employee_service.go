package service

import (
	"errors"
	"fmt"
	"u_kom_be/internal/converter"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
	"u_kom_be/internal/utils"
)

type EmployeeService interface {
	CreateEmployee(req request.EmployeeCreateRequest) (*response.EmployeeDetailResponse, error)
	GetEmployeeByID(id string) (*response.EmployeeDetailResponse, error)
	GetAllEmployees(search string) ([]response.EmployeeListResponse, error)
	UpdateEmployee(id string, req request.EmployeeUpdateRequest) (*response.EmployeeDetailResponse, error)
	DeleteEmployee(id string) error

	// Method untuk "Project A" (Menautkan Akun)
	LinkUser(employeeID string, userID string) error
	UnlinkUser(employeeID string) error
}

type employeeService struct {
	employeeRepo repository.EmployeeRepository
	userRepo     repository.UserRepository // Dependensi untuk validasi user
	encryptUtil  utils.EncryptionUtil
	converter    converter.EmployeeConverterInterface
}

func NewEmployeeService(
	employeeRepo repository.EmployeeRepository,
	userRepo repository.UserRepository, // Tambahkan parameter
	encryptUtil utils.EncryptionUtil,
	converter converter.EmployeeConverterInterface,
) EmployeeService {
	return &employeeService{
		employeeRepo: employeeRepo,
		userRepo:     userRepo, // Inject dependensi
		encryptUtil:  encryptUtil,
		converter:    converter,
	}
}

// CreateEmployee menangani pembuatan pegawai baru
func (s *employeeService) CreateEmployee(req request.EmployeeCreateRequest) (*response.EmployeeDetailResponse, error) {
	// 1. Validasi Duplikat (NIP & Phone)
	if req.NIP != nil && *req.NIP != "" {
		if existing, _ := s.employeeRepo.FindByNIP(*req.NIP); existing != nil {
			return nil, errors.New("nip already exists")
		}
	}
	if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		if existing, _ := s.employeeRepo.FindByPhone(*req.PhoneNumber); existing != nil {
			return nil, errors.New("phone number already exists")
		}
	}

	// 2. Enkripsi & Hash NIK
	encryptedNIK := ""
	nikHash := ""
	if req.NIK != "" {
		// a. Cek Unik via Blind Index
		var err error
		nikHash, err = s.encryptUtil.Hash(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to hash nik: %w", err)
		}
		existing, err := s.employeeRepo.FindByNIKHash(nikHash)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("nik already exists")
		}

		// b. Enkripsi
		encryptedNIK, err = s.encryptUtil.Encrypt(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt nik: %w", err)
		}
	}

	// 3. Buat Domain Object (Mapping DTO ke Domain)
	employee := &domain.Employee{
		FullName:         req.FullName,
		NIP:              req.NIP,
		JobTitle:         req.JobTitle,
		NIK:              encryptedNIK, // <-- Simpan data terenkripsi
		NIKHash:          nikHash,      // <-- Simpan hash untuk validasi
		Gender:           req.Gender,
		PhoneNumber:      req.PhoneNumber,
		Address:          req.Address,
		DateOfBirth:      req.DateOfBirth,
		JoinDate:         req.JoinDate,
		EmploymentStatus: req.EmploymentStatus,
		// UserID sengaja NULL saat dibuat
	}

	// 4. Panggil Repository
	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, err
	}

	// 5. Ambil data yang baru dibuat
	createdEmployee, err := s.employeeRepo.FindByID(employee.ID)
	if err != nil {
		return nil, err
	}
	if createdEmployee == nil {
		return nil, errors.New("failed to retrieve created employee")
	}

	// 6. Konversi ke Response Detail
	resp := s.converter.ToEmployeeDetailResponse(createdEmployee)

	// === LOGIC MAPPING USER ===
	if employee.User.ID != "" {
		resp.User = &response.UserLinkedResponse{
			ID:       employee.User.ID,
			Username: employee.User.Username,
			Name:     employee.User.Name,
			Email:    employee.User.Email,
		}
	} else {
		resp.User = nil
	}

	return resp, nil
}

// GetEmployeeByID mengambil satu data pegawai
func (s *employeeService) GetEmployeeByID(id string) (*response.EmployeeDetailResponse, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, errors.New("employee not found")
	}
	// Panggil konverter (akan mendekripsi NIK)
	resp := s.converter.ToEmployeeDetailResponse(employee)

	// === LOGIC MAPPING USER ===
	if employee.User.ID != "" {
		resp.User = &response.UserLinkedResponse{
			ID:       employee.User.ID,
			Username: employee.User.Username,
			Name:     employee.User.Name,
			Email:    employee.User.Email,
		}
	} else {
		resp.User = nil
	}

	return resp, nil
}

// GetAllEmployees mengambil semua data pegawai (ringkas)
func (s *employeeService) GetAllEmployees(search string) ([]response.EmployeeListResponse, error) {
	employees, err := s.employeeRepo.FindAll(search)
	if err != nil {
		return nil, err
	}
	// Panggil konverter untuk list (tanpa NIK)
	return s.converter.ToEmployeeListResponses(employees), nil
}

// UpdateEmployee memperbarui data pegawai
func (s *employeeService) UpdateEmployee(id string, req request.EmployeeUpdateRequest) (*response.EmployeeDetailResponse, error) {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, errors.New("employee not found")
	}

	// Update fields jika disediakan
	if req.FullName != "" {
		employee.FullName = req.FullName
	}
	if req.JobTitle != "" {
		employee.JobTitle = req.JobTitle
	}

	// Validasi duplikat NIP (jika NIP diubah)
	if req.NIP != nil {
		if (employee.NIP == nil || *req.NIP != *employee.NIP) && *req.NIP != "" {
			if existing, _ := s.employeeRepo.FindByNIP(*req.NIP); existing != nil {
				return nil, errors.New("nip already exists")
			}
		}
		employee.NIP = req.NIP
	}

	// Validasi duplikat Telepon (jika diubah)
	if req.PhoneNumber != nil {
		if (employee.PhoneNumber == nil || *req.PhoneNumber != *employee.PhoneNumber) && *req.PhoneNumber != "" {
			if existing, _ := s.employeeRepo.FindByPhone(*req.PhoneNumber); existing != nil {
				return nil, errors.New("phone number already exists")
			}
		}
		employee.PhoneNumber = req.PhoneNumber
	}

	// Enkripsi & Hash NIK jika diperbarui
	if req.NIK != "" {
		// Validasi unik
		nikHash, err := s.encryptUtil.Hash(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to hash nik: %w", err)
		}

		// Cek apakah hash sudah ada di record LAIN
		existing, err := s.employeeRepo.FindByNIKHash(nikHash)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("nik already exists")
		}

		encryptedNIK, err := s.encryptUtil.Encrypt(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt nik: %w", err)
		}
		employee.NIK = encryptedNIK
		employee.NIKHash = nikHash
	}

	// Update field lainnya
	if req.Gender != "" {
		employee.Gender = req.Gender
	}
	if req.Address != "" {
		employee.Address = req.Address
	}
	if req.EmploymentStatus != "" {
		employee.EmploymentStatus = req.EmploymentStatus
	}
	if req.DateOfBirth != nil {
		employee.DateOfBirth = req.DateOfBirth
	}
	if req.JoinDate != nil {
		employee.JoinDate = req.JoinDate
	}

	// Simpan perubahan
	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, err
	}

	// Ambil data yang sudah diupdate
	updatedEmployee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := s.converter.ToEmployeeDetailResponse(updatedEmployee)

	// === LOGIC MAPPING USER ===
	if employee.User.ID != "" {
		resp.User = &response.UserLinkedResponse{
			ID:       employee.User.ID,
			Username: employee.User.Username,
			Name:     employee.User.Name,
			Email:    employee.User.Email,
		}
	} else {
		resp.User = nil
	}

	return resp, nil
}

// DeleteEmployee menghapus data pegawai
func (s *employeeService) DeleteEmployee(id string) error {
	employee, err := s.employeeRepo.FindByID(id)
	if err != nil {
		return err
	}
	if employee == nil {
		return errors.New("employee not found")
	}

	// TODO: Tambahkan logika bisnis,
	// misal: "tidak bisa hapus pegawai jika dia adalah Wali Kelas aktif"

	return s.employeeRepo.Delete(id)
}

// --- Method untuk "Project A" (Menautkan Akun) ---

// LinkUser menautkan profil Employee ke akun User
func (s *employeeService) LinkUser(employeeID string, userID string) error {
	// 1. Cek apakah Employee ada
	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return err
	}
	if employee == nil {
		return errors.New("employee not found")
	}

	// 2. Cek apakah User ada
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 3. Cek apakah User tersebut sudah ditautkan ke Employee LAIN
	existingEmployee, _ := s.employeeRepo.FindByUserID(userID)
	if existingEmployee != nil && existingEmployee.ID != employeeID {
		return errors.New("this user account is already linked to another employee")
	}

	// 4. Tautkan akun
	return s.employeeRepo.SetUserID(employeeID, &userID)
}

// UnlinkUser menghapus tautan Employee dari akun User
func (s *employeeService) UnlinkUser(employeeID string) error {
	// 1. Cek apakah Employee ada
	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return err
	}
	if employee == nil {
		return errors.New("employee not found")
	}

	// 2. Hapus tautan (set user_id ke NULL)
	return s.employeeRepo.SetUserID(employeeID, nil)
}
