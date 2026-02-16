package service

import (
	"fmt"
	"smart_school_be/internal/apperrors"
	"smart_school_be/internal/converter"
	"smart_school_be/internal/model/domain"
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/model/response"
	"smart_school_be/internal/repository"
	"smart_school_be/internal/utils"
	"strings"
)

type ParentService interface {
	CreateParent(req request.ParentCreateRequest) (*response.ParentDetailResponse, error)
	GetParentByID(id string) (*response.ParentDetailResponse, error)
	GetAllParents(search string, pagination request.PaginationRequest) (*response.PaginatedData, error)
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
	// 1. Validasi Duplikat (Phone & Email) - jika ada
	if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		if existing, _ := s.parentRepo.FindByPhone(*req.PhoneNumber); existing != nil {
			return nil, apperrors.NewConflictError("Phone number already exists")
		}
	}
	if req.Email != nil && *req.Email != "" {
		if existing, _ := s.parentRepo.FindByEmail(*req.Email); existing != nil {
			return nil, apperrors.NewConflictError("Email already exists")
		}
	}

	// 2. Enkripsi NIK & Hash
	var encryptedNIK *string
	var nikHash *string

	if req.NIK != nil && *req.NIK != "" {
		// a. Hash & Check Unique
		hash, err := s.encryptionUtil.Hash(*req.NIK)
		if err != nil {
			return nil, fmt.Errorf("Failed to hash NIK: %w", err)
		}

		existing, err := s.parentRepo.FindByNIKHash(hash)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, apperrors.NewConflictError("NIK already exists")
		}
		nikHash = &hash

		// b. Encrypt
		encrypted, err := s.encryptionUtil.Encrypt(*req.NIK)
		if err != nil {
			return nil, fmt.Errorf("Failed to encrypt NIK: %w", err)
		}
		encryptedNIK = &encrypted
	}

	// 3. Buat Domain Object
	// Helper untuk set default LifeStatus
	lifeStatus := "alive"
	if req.LifeStatus != nil {
		lifeStatus = *req.LifeStatus
	}

	parent := &domain.Parent{
		FullName:       req.FullName,
		NIK:            encryptedNIK, // <-- Simpan data terenkripsi (pointer)
		NIKHash:        nikHash,      // <-- Simpan hash (pointer)
		Gender:         req.Gender,
		PlaceOfBirth:   req.PlaceOfBirth,
		DateOfBirth:    req.DateOfBirth,
		LifeStatus:     &lifeStatus,
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
		return nil, apperrors.NewInternalError("Failed to retrieve created parent")
	}

	// 6. Konversi ke Response Detail
	resp := s.converter.ToParentDetailResponse(createdParent)

	// Cek apakah parent punya user_id (terhubung ke akun)
	if parent.User.ID != "" {
		resp.User = &response.UserLinkedResponse{
			ID:       parent.User.ID,
			Username: parent.User.Username,
			Name:     parent.User.Name,
			Email:    parent.User.Email,
		}
	} else {
		resp.User = nil
	}

	return resp, nil
}

// GetParentByID mengambil satu data orang tua
func (s *parentService) GetParentByID(id string) (*response.ParentDetailResponse, error) {
	// 1. Repo sudah Preload User
	parent, err := s.parentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, apperrors.NewNotFoundError("Parent not found")
	}

	// 2. Konversi Parent dasar
	resp := s.converter.ToParentDetailResponse(parent)

	// 3. PERBAIKAN: Mapping Data User (Manual)
	// Cek apakah parent punya user_id (terhubung ke akun)
	if parent.User.ID != "" {
		resp.User = &response.UserLinkedResponse{
			ID:       parent.User.ID,
			Username: parent.User.Username,
			Name:     parent.User.Name,
			Email:    parent.User.Email,
		}
	} else {
		resp.User = nil
	}

	return resp, nil
}

// GetAllParents mengambil semua data orang tua (ringkas) dengan pagination
func (s *parentService) GetAllParents(search string, pagination request.PaginationRequest) (*response.PaginatedData, error) {
	limit := pagination.GetLimit()
	offset := pagination.GetOffset()

	parents, total, err := s.parentRepo.FindAll(search, limit, offset)
	if err != nil {
		return nil, err
	}
	// Panggil konverter untuk response list
	data := s.converter.ToParentListResponses(parents)
	paginatedData := response.NewPaginatedData(data, total, pagination.GetPage(), limit)
	return &paginatedData, nil
}

// UpdateParent memperbarui data orang tua
func (s *parentService) UpdateParent(id string, req request.ParentUpdateRequest) (*response.ParentDetailResponse, error) {
	parent, err := s.parentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, apperrors.NewNotFoundError("Parent not found")
	}

	// Update fields jika disediakan
	if req.FullName != "" {
		parent.FullName = req.FullName
	}

	// Validasi duplikat untuk PhoneNumber
	// Jika nil atau empty string dari request, kita set ke nil (hapus phone number)
	// Jika ada value, validasi uniqueness terlebih dahulu
	if req.PhoneNumber != nil && *req.PhoneNumber != "" {
		// Ada value baru, cek duplikat
		if parent.PhoneNumber == nil || *req.PhoneNumber != *parent.PhoneNumber {
			if existing, _ := s.parentRepo.FindByPhone(*req.PhoneNumber); existing != nil {
				return nil, apperrors.NewConflictError("Phone number already exists")
			}
		}
	}
	// Update phone number (termasuk jika nil atau empty untuk set null)
	parent.PhoneNumber = req.PhoneNumber

	// Validasi duplikat untuk Email
	if req.Email != nil && *req.Email != "" {
		// Ada value baru, cek duplikat
		if parent.Email == nil || *req.Email != *parent.Email {
			if existing, _ := s.parentRepo.FindByEmail(*req.Email); existing != nil {
				return nil, apperrors.NewConflictError("Email already exists")
			}
		}
	}
	// Update email (termasuk jika nil atau empty untuk set null)
	parent.Email = req.Email

	// Enkripsi NIK jika diperbarui
	if req.NIK != nil {
		// Jika kosong string, kita anggap hapus NIK? Atau validasi di FE?
		// Asumsi: jika pointer != nil dan string != "", kita update. Jika "", kita bisa set nil atau biarkan kosong.
		// Sesuai logic Create, NIK boleh kosong/nil.
		if *req.NIK != "" {
			// Hitung hash baru
			newHash, err := s.encryptionUtil.Hash(*req.NIK)
			if err != nil {
				return nil, fmt.Errorf("Failed to hash NIK: %w", err)
			}

			// Cek keunikan jika hash berubah atau sebelumnya null
			isDifferent := parent.NIKHash == nil || *parent.NIKHash != newHash
			if isDifferent {
				existing, err := s.parentRepo.FindByNIKHash(newHash)
				if err != nil {
					return nil, err
				}
				if existing != nil && existing.ID != id {
					return nil, apperrors.NewConflictError("NIK already exists")
				}
			}

			parent.NIKHash = &newHash

			// Enkripsi
			encryptedNIK, err := s.encryptionUtil.Encrypt(*req.NIK)
			if err != nil {
				return nil, fmt.Errorf("Failed to encrypt NIK: %w", err)
			}
			parent.NIK = &encryptedNIK
		} else {
			// Jika explicit empty string, mungkin user ingin menghapus NIK?
			parent.NIK = nil
			parent.NIKHash = nil
		}
	}

	// Update field lainnya - langsung assign untuk bisa null dari JSON
	parent.Gender = req.Gender
	parent.PlaceOfBirth = req.PlaceOfBirth
	parent.DateOfBirth = req.DateOfBirth
	parent.LifeStatus = req.LifeStatus
	parent.MaritalStatus = req.MaritalStatus
	parent.EducationLevel = req.EducationLevel
	parent.Occupation = req.Occupation
	parent.IncomeRange = req.IncomeRange
	parent.Address = req.Address
	parent.RT = req.RT
	parent.RW = req.RW
	parent.SubDistrict = req.SubDistrict
	parent.District = req.District
	parent.City = req.City
	parent.Province = req.Province
	parent.PostalCode = req.PostalCode

	// Simpan perubahan
	if err := s.parentRepo.Update(parent); err != nil {
		return nil, err
	}

	// Ambil data yang sudah diupdate
	updatedParent, err := s.parentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resp := s.converter.ToParentDetailResponse(updatedParent)

	// Cek apakah parent punya user_id (terhubung ke akun)
	if parent.User.ID != "" {
		resp.User = &response.UserLinkedResponse{
			ID:       parent.User.ID,
			Username: parent.User.Username,
			Name:     parent.User.Name,
			Email:    parent.User.Email,
		}
	} else {
		resp.User = nil
	}

	return resp, nil
}

// DeleteParent menghapus data orang tua
func (s *parentService) DeleteParent(id string) error {
	parent, err := s.parentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if parent == nil {
		return apperrors.NewNotFoundError("Parent not found")
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
		return apperrors.NewNotFoundError("Parent not found")
	}

	// 2. Cek apakah User ada
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.NewNotFoundError("User not found")
	}

	// 3. Tautkan akun (Kita andalkan UNIQUE constraint di DB untuk error duplikat)
	if err := s.parentRepo.SetUserID(parentID, &userID); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return apperrors.NewConflictError("This user account is already linked to another parent")
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
		return apperrors.NewNotFoundError("Parent not found")
	}

	// 2. Hapus tautan (set user_id ke NULL)
	return s.parentRepo.SetUserID(parentID, nil)
}
