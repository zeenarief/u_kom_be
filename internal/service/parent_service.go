package service

import (
	"belajar-golang/internal/converter"
	"belajar-golang/internal/model/domain"
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/model/response"
	"belajar-golang/internal/repository"
	"belajar-golang/internal/utils"
	"errors"
	"fmt"
	"strings"
)

type ParentService interface {
	CreateParent(req request.ParentCreateRequest) (*response.ParentDetailResponse, error)
	GetParentByID(id string) (*response.ParentDetailResponse, error)
	GetAllParents() ([]response.ParentListResponse, error)
	UpdateParent(id string, req request.ParentUpdateRequest) (*response.ParentDetailResponse, error)
	DeleteParent(id string) error
	LinkUser(parentID string, userID string) error
	UnlinkUser(parentID string) error
}

type parentService struct {
	parentRepo     repository.ParentRepository
	userRepo       repository.UserRepository
	encryptionUtil utils.EncryptionUtil
	converter      converter.ParentConverterInterface
}

func NewParentService(
	parentRepo repository.ParentRepository,
	userRepo repository.UserRepository,
	encryptionUtil utils.EncryptionUtil,
	converter converter.ParentConverterInterface,
) ParentService {
	return &parentService{
		parentRepo:     parentRepo,
		userRepo:       userRepo,
		encryptionUtil: encryptionUtil,
		converter:      converter,
	}
}

// CreateParent menangani pembuatan data orang tua baru
func (s *parentService) CreateParent(req request.ParentCreateRequest) (*response.ParentDetailResponse, error) {
	// 1. Validasi Duplikat (Phone & Email)
	if req.PhoneNumber != "" {
		if existing, _ := s.parentRepo.FindByPhone(req.PhoneNumber); existing != nil {
			return nil, errors.New("phone number already exists")
		}
	}
	if req.Email != "" {
		if existing, _ := s.parentRepo.FindByEmail(req.Email); existing != nil {
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
	parent := &domain.Parent{
		FullName:       req.FullName,
		NIK:            encryptedNIK, // <-- Simpan data terenkripsi
		Gender:         req.Gender,
		PlaceOfBirth:   req.PlaceOfBirth,
		DateOfBirth:    req.DateOfBirth,
		LifeStatus:     req.LifeStatus,
		MaritalStatus:  req.MaritalStatus,
		PhoneNumber:    req.PhoneNumber,
		Email:          req.Email,
		EducationLevel: req.EducationLevel,
		Occupation:     req.Occupation,
		IncomeRange:    req.IncomeRange,
		Address:        req.Address,
		RT:             req.RT,
		RW:             req.RW,
		SubDistrict:    req.SubDistrict,
		District:       req.District,
		City:           req.City,
		Province:       req.Province,
		PostalCode:     req.PostalCode,
	}

	// Menetapkan default 'alive' jika tidak diset, sesuai skema DB
	if parent.LifeStatus == "" {
		parent.LifeStatus = "alive"
	}

	// 4. Panggil Repository
	if err := s.parentRepo.Create(parent); err != nil {
		return nil, err
	}

	// 5. Ambil data yang baru dibuat
	createdParent, err := s.parentRepo.FindByID(parent.ID)
	if err != nil {
		return nil, err
	}
	if createdParent == nil {
		return nil, errors.New("failed to retrieve created parent")
	}

	// 6. Konversi ke Response Detail
	return s.converter.ToParentDetailResponse(createdParent), nil
}

// GetParentByID mengambil satu data orang tua
func (s *parentService) GetParentByID(id string) (*response.ParentDetailResponse, error) {
	parent, err := s.parentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, errors.New("parent not found")
	}
	// Panggil konverter untuk response detail (dengan dekripsi NIK)
	return s.converter.ToParentDetailResponse(parent), nil
}

// GetAllParents mengambil semua data orang tua (ringkas)
func (s *parentService) GetAllParents() ([]response.ParentListResponse, error) {
	parents, err := s.parentRepo.FindAll()
	if err != nil {
		return nil, err
	}
	// Panggil konverter untuk response list (ringkas, tanpa NIK)
	return s.converter.ToParentListResponses(parents), nil
}

// UpdateParent memperbarui data orang tua
func (s *parentService) UpdateParent(id string, req request.ParentUpdateRequest) (*response.ParentDetailResponse, error) {
	parent, err := s.parentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, errors.New("parent not found")
	}

	// Update fields jika disediakan
	if req.FullName != "" {
		parent.FullName = req.FullName
	}

	// Validasi duplikat baru
	if req.PhoneNumber != "" && req.PhoneNumber != parent.PhoneNumber {
		if existing, _ := s.parentRepo.FindByPhone(req.PhoneNumber); existing != nil {
			return nil, errors.New("phone number already exists")
		}
		parent.PhoneNumber = req.PhoneNumber
	}
	if req.Email != "" && req.Email != parent.Email {
		if existing, _ := s.parentRepo.FindByEmail(req.Email); existing != nil {
			return nil, errors.New("email already exists")
		}
		parent.Email = req.Email
	}

	// Enkripsi NIK jika diperbarui
	if req.NIK != "" {
		encryptedNIK, err := s.encryptionUtil.Encrypt(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt nik: %w", err)
		}
		parent.NIK = encryptedNIK
	}

	// Update field lainnya (pola `!= ""`)
	if req.Gender != "" {
		parent.Gender = req.Gender
	}
	if req.PlaceOfBirth != "" {
		parent.PlaceOfBirth = req.PlaceOfBirth
	}
	if !req.DateOfBirth.IsZero() {
		parent.DateOfBirth = req.DateOfBirth
	}
	if req.LifeStatus != "" {
		parent.LifeStatus = req.LifeStatus
	}
	if req.MaritalStatus != "" {
		parent.MaritalStatus = req.MaritalStatus
	}
	if req.EducationLevel != "" {
		parent.EducationLevel = req.EducationLevel
	}
	if req.Occupation != "" {
		parent.Occupation = req.Occupation
	}
	if req.IncomeRange != "" {
		parent.IncomeRange = req.IncomeRange
	}
	if req.Address != "" {
		parent.Address = req.Address
	}
	if req.RT != "" {
		parent.RT = req.RT
	}
	if req.RW != "" {
		parent.RW = req.RW
	}
	if req.SubDistrict != "" {
		parent.SubDistrict = req.SubDistrict
	}
	if req.District != "" {
		parent.District = req.District
	}
	if req.City != "" {
		parent.City = req.City
	}
	if req.Province != "" {
		parent.Province = req.Province
	}
	if req.PostalCode != "" {
		parent.PostalCode = req.PostalCode
	}

	// Simpan perubahan
	if err := s.parentRepo.Update(parent); err != nil {
		return nil, err
	}

	// Ambil data yang sudah diupdate
	updatedParent, err := s.parentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.converter.ToParentDetailResponse(updatedParent), nil
}

// DeleteParent menghapus data orang tua
func (s *parentService) DeleteParent(id string) error {
	parent, err := s.parentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if parent == nil {
		return errors.New("parent not found")
	}

	return s.parentRepo.Delete(id)
}

// LinkUser menautkan profil Parent ke akun User
func (s *parentService) LinkUser(parentID string, userID string) error {
	// 1. Cek apakah Parent ada
	parent, err := s.parentRepo.FindByID(parentID)
	if err != nil {
		return err
	}
	if parent == nil {
		return errors.New("parent not found")
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
	if err := s.parentRepo.SetUserID(parentID, &userID); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.New("this user account is already linked to another parent")
		}
		return err
	}
	return nil
}

// UnlinkUser menghapus tautan Parent dari akun User
func (s *parentService) UnlinkUser(parentID string) error {
	// 1. Cek apakah Parent ada
	parent, err := s.parentRepo.FindByID(parentID)
	if err != nil {
		return err
	}
	if parent == nil {
		return errors.New("parent not found")
	}

	// 2. Hapus tautan (set user_id ke NULL)
	return s.parentRepo.SetUserID(parentID, nil)
}
