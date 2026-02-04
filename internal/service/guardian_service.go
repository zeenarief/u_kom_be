package service

import (
	"errors"
	"fmt"
	"strings"
	"u_kom_be/internal/converter"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
	"u_kom_be/internal/utils"
)

type GuardianService interface {
	CreateGuardian(req request.GuardianCreateRequest) (*response.GuardianDetailResponse, error)
	GetGuardianByID(id string) (*response.GuardianDetailResponse, error)
	GetAllGuardians(search string) ([]response.GuardianListResponse, error)
	UpdateGuardian(id string, req request.GuardianUpdateRequest) (*response.GuardianDetailResponse, error)
	DeleteGuardian(id string) error
	LinkUser(guardianID string, userID string) error
	UnlinkUser(guardianID string) error
}

type guardianService struct {
	guardianRepo   repository.GuardianRepository
	userRepo       repository.UserRepository
	encryptionUtil utils.EncryptionUtil
	converter      converter.GuardianConverterInterface
}

func NewGuardianService(
	guardianRepo repository.GuardianRepository,
	userRepo repository.UserRepository,
	encryptionUtil utils.EncryptionUtil,
	converter converter.GuardianConverterInterface,
) GuardianService {
	return &guardianService{
		guardianRepo:   guardianRepo,
		userRepo:       userRepo,
		encryptionUtil: encryptionUtil,
		converter:      converter,
	}
}

// CreateGuardian menangani pembuatan data wali baru
func (s *guardianService) CreateGuardian(req request.GuardianCreateRequest) (*response.GuardianDetailResponse, error) {
	// 1. Validasi Duplikat (Phone & Email)
	if req.PhoneNumber != "" { // Phone number wajib ada
		if existing, _ := s.guardianRepo.FindByPhone(req.PhoneNumber); existing != nil {
			return nil, errors.New("phone number already exists")
		}
	}
	if req.Email != "" {
		if existing, _ := s.guardianRepo.FindByEmail(req.Email); existing != nil {
			return nil, errors.New("email already exists")
		}
	}

	// 2. Enkripsi NIK
	encryptedNIK := ""
	if req.NIK != "" {
		var err error
		encryptedNIK, err = s.encryptionUtil.Encrypt(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt nik: %w", err)
		}
	}

	// 3. Buat Domain Object
	guardian := &domain.Guardian{
		FullName:              req.FullName,
		NIK:                   encryptedNIK, // <-- Simpan data terenkripsi
		Gender:                req.Gender,
		PhoneNumber:           req.PhoneNumber,
		Email:                 req.Email,
		Address:               req.Address,
		RT:                    req.RT,
		RW:                    req.RW,
		SubDistrict:           req.SubDistrict,
		District:              req.District,
		City:                  req.City,
		Province:              req.Province,
		PostalCode:            req.PostalCode,
		RelationshipToStudent: req.RelationshipToStudent,
	}

	// 4. Panggil Repository
	if err := s.guardianRepo.Create(guardian); err != nil {
		return nil, err
	}

	// 5. Ambil data yang baru dibuat
	createdGuardian, err := s.guardianRepo.FindByID(guardian.ID)
	if err != nil {
		return nil, err
	}
	if createdGuardian == nil {
		return nil, errors.New("failed to retrieve created guardian")
	}

	// 6. Konversi ke Response Detail
	resp := s.converter.ToGuardianDetailResponse(createdGuardian)

	// MAPPING USER
	if guardian.User.ID != "" {
		resp.User = &response.UserLinkedResponse{
			ID:       guardian.User.ID,
			Username: guardian.User.Username,
			Name:     guardian.User.Name,
			Email:    guardian.User.Email,
		}
	} else {
		resp.User = nil
	}

	return resp, nil
}

// GetGuardianByID mengambil satu data wali
func (s *guardianService) GetGuardianByID(id string) (*response.GuardianDetailResponse, error) {
	guardian, err := s.guardianRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if guardian == nil {
		return nil, errors.New("guardian not found")
	}
	// Panggil konverter untuk response detail (dengan dekripsi NIK)
	resp := s.converter.ToGuardianDetailResponse(guardian)

	// MAPPING USER
	if guardian.User.ID != "" {
		resp.User = &response.UserLinkedResponse{
			ID:       guardian.User.ID,
			Username: guardian.User.Username,
			Name:     guardian.User.Name,
			Email:    guardian.User.Email,
		}
	} else {
		resp.User = nil
	}

	return resp, nil
}

// GetAllGuardians mengambil semua data wali (ringkas)
func (s *guardianService) GetAllGuardians(search string) ([]response.GuardianListResponse, error) {
	guardians, err := s.guardianRepo.FindAll(search)
	if err != nil {
		return nil, err
	}
	// Panggil konverter untuk response list (ringkas, tanpa NIK)
	return s.converter.ToGuardianListResponses(guardians), nil
}

// UpdateGuardian memperbarui data wali
func (s *guardianService) UpdateGuardian(id string, req request.GuardianUpdateRequest) (*response.GuardianDetailResponse, error) {
	guardian, err := s.guardianRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if guardian == nil {
		return nil, errors.New("guardian not found")
	}

	// Update fields jika disediakan
	if req.FullName != "" {
		guardian.FullName = req.FullName
	}

	// Validasi duplikat baru
	if req.PhoneNumber != "" && req.PhoneNumber != guardian.PhoneNumber {
		if existing, _ := s.guardianRepo.FindByPhone(req.PhoneNumber); existing != nil {
			return nil, errors.New("phone number already exists")
		}
		guardian.PhoneNumber = req.PhoneNumber
	}
	if req.Email != "" && req.Email != guardian.Email {
		if existing, _ := s.guardianRepo.FindByEmail(req.Email); existing != nil {
			return nil, errors.New("email already exists")
		}
		guardian.Email = req.Email
	}

	// Enkripsi NIK jika diperbarui
	if req.NIK != "" {
		encryptedNIK, err := s.encryptionUtil.Encrypt(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt nik: %w", err)
		}
		guardian.NIK = encryptedNIK
	}

	// Update field lainnya (pola `!= ""`)
	if req.Gender != "" {
		guardian.Gender = req.Gender
	}
	if req.Address != "" {
		guardian.Address = req.Address
	}
	if req.RT != "" {
		guardian.RT = req.RT
	}
	if req.RW != "" {
		guardian.RW = req.RW
	}
	if req.SubDistrict != "" {
		guardian.SubDistrict = req.SubDistrict
	}
	if req.District != "" {
		guardian.District = req.District
	}
	if req.City != "" {
		guardian.City = req.City
	}
	if req.Province != "" {
		guardian.Province = req.Province
	}
	if req.PostalCode != "" {
		guardian.PostalCode = req.PostalCode
	}
	if req.RelationshipToStudent != "" {
		guardian.RelationshipToStudent = req.RelationshipToStudent
	}

	// Simpan perubahan
	if err := s.guardianRepo.Update(guardian); err != nil {
		return nil, err
	}

	// Ambil data yang sudah diupdate
	updatedGuardian, err := s.guardianRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := s.converter.ToGuardianDetailResponse(updatedGuardian)

	// MAPPING USER
	if guardian.User.ID != "" {
		resp.User = &response.UserLinkedResponse{
			ID:       guardian.User.ID,
			Username: guardian.User.Username,
			Name:     guardian.User.Name,
			Email:    guardian.User.Email,
		}
	} else {
		resp.User = nil
	}

	return resp, nil
}

// DeleteGuardian menghapus data wali
func (s *guardianService) DeleteGuardian(id string) error {
	guardian, err := s.guardianRepo.FindByID(id)
	if err != nil {
		return err
	}
	if guardian == nil {
		return errors.New("guardian not found")
	}

	return s.guardianRepo.Delete(id)
}

// LinkUser menautkan profil Guardian ke akun User
func (s *guardianService) LinkUser(guardianID string, userID string) error {
	// 1. Cek apakah Guardian ada
	guardian, err := s.guardianRepo.FindByID(guardianID)
	if err != nil {
		return err
	}
	if guardian == nil {
		return errors.New("guardian not found")
	}

	// 2. Cek apakah User ada
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 3. Tautkan akun (Kita andalkan UNIQUE constraint di DB untuk error duplikat)
	if err := s.guardianRepo.SetUserID(guardianID, &userID); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.New("this user account is already linked to another guardian")
		}
		return err
	}
	return nil
}

// UnlinkUser menghapus tautan Guardian dari akun User
func (s *guardianService) UnlinkUser(guardianID string) error {
	// 1. Cek apakah Guardian ada
	guardian, err := s.guardianRepo.FindByID(guardianID)
	if err != nil {
		return err
	}
	if guardian == nil {
		return errors.New("guardian not found")
	}

	// 2. Hapus tautan (set user_id ke NULL)
	return s.guardianRepo.SetUserID(guardianID, nil)
}
