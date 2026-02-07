package service

import (
	"errors"
	"fmt"
	"strings"
	"u_kom_be/internal/apperrors"
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
	if req.PhoneNumber != nil && *req.PhoneNumber != "" { // Phone number optional
		if existing, _ := s.guardianRepo.FindByPhone(*req.PhoneNumber); existing != nil {
			return nil, apperrors.NewConflictError("Phone number already exists")
		}
	}
	if req.Email != nil && *req.Email != "" {
		if existing, _ := s.guardianRepo.FindByEmail(*req.Email); existing != nil {
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

		existing, err := s.guardianRepo.FindByNIKHash(hash)
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
	guardian := &domain.Guardian{
		FullName:              req.FullName,
		NIK:                   encryptedNIK, // <-- Simpan data terenkripsi
		NIKHash:               nikHash,      // <-- Simpan hash
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
		return nil, apperrors.NewInternalError("Failed to retrieve created guardian")
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
		return nil, apperrors.NewNotFoundError("Guardian not found")
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
		return nil, apperrors.NewNotFoundError("Guardian not found")
	}

	// Update fields jika disediakan
	if req.FullName != "" {
		guardian.FullName = req.FullName
	}

	// Validasi duplikat baru
	if req.PhoneNumber != nil {
		if guardian.PhoneNumber == nil || *req.PhoneNumber != *guardian.PhoneNumber {
			if *req.PhoneNumber != "" {
				if existing, _ := s.guardianRepo.FindByPhone(*req.PhoneNumber); existing != nil {
					return nil, apperrors.NewConflictError("Phone number already exists")
				}
			}
			guardian.PhoneNumber = req.PhoneNumber
		}
	}

	if req.Email != nil {
		if guardian.Email == nil || *req.Email != *guardian.Email {
			if *req.Email != "" {
				if existing, _ := s.guardianRepo.FindByEmail(*req.Email); existing != nil {
					return nil, apperrors.NewConflictError("Email already exists")
				}
			}
			guardian.Email = req.Email
		}
	}

	// Enkripsi NIK jika diperbarui
	if req.NIK != nil {
		if *req.NIK != "" {
			// Hitung hash baru
			newHash, err := s.encryptionUtil.Hash(*req.NIK)
			if err != nil {
				return nil, fmt.Errorf("Failed to hash NIK: %w", err)
			}

			// Cek keunikan jika hash berubah atau sebelumnya null
			isDifferent := guardian.NIKHash == nil || *guardian.NIKHash != newHash
			if isDifferent {
				existing, err := s.guardianRepo.FindByNIKHash(newHash)
				if err != nil {
					return nil, err
				}
				if existing != nil && existing.ID != id {
					return nil, apperrors.NewConflictError("NIK already exists")
				}
			}

			guardian.NIKHash = &newHash

			// Enkripsi
			encryptedNIK, err := s.encryptionUtil.Encrypt(*req.NIK)
			if err != nil {
				return nil, fmt.Errorf("Failed to encrypt NIK: %w", err)
			}
			guardian.NIK = &encryptedNIK
		} else {
			// Jika explicit empty string, hapus NIK
			guardian.NIK = nil
			guardian.NIKHash = nil
		}
	}

	// Update field lainnya (pola `!= nil`)
	if req.Gender != nil {
		guardian.Gender = req.Gender
	}
	if req.Address != nil {
		guardian.Address = req.Address
	}
	if req.RT != nil {
		guardian.RT = req.RT
	}
	if req.RW != nil {
		guardian.RW = req.RW
	}
	if req.SubDistrict != nil {
		guardian.SubDistrict = req.SubDistrict
	}
	if req.District != nil {
		guardian.District = req.District
	}
	if req.City != nil {
		guardian.City = req.City
	}
	if req.Province != nil {
		guardian.Province = req.Province
	}
	if req.PostalCode != nil {
		guardian.PostalCode = req.PostalCode
	}
	if req.RelationshipToStudent != nil {
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
		return apperrors.NewNotFoundError("Guardian not found")
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
		return apperrors.NewNotFoundError("Guardian not found")
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
	if err := s.guardianRepo.SetUserID(guardianID, &userID); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return apperrors.NewConflictError("This user account is already linked to another guardian")
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
		return errors.New("Guardian not found")
	}

	// 2. Hapus tautan (set user_id ke NULL)
	return s.guardianRepo.SetUserID(guardianID, nil)
}
