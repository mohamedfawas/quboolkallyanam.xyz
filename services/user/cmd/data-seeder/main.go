package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/utils/validation"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/config"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/user/internal/domain/entity"
	"gorm.io/gorm"
)

const (
	BATCH_SIZE      = 1000   // Insert 1000 records per batch
	DEFAULT_RECORDS = 100000 // Default: 1 lakh records
)

func main() {
	// Get number of records from command line argument
	targetRecords := DEFAULT_RECORDS
	if len(os.Args) > 1 {
		if records, err := strconv.Atoi(os.Args[1]); err == nil {
			targetRecords = records
		}
	}

	fmt.Printf("🚀 Starting bulk data insertion for %d user profiles...\n", targetRecords)

	// Initialize database connection
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Generate and insert data
	start := time.Now()
	err = generateBulkUserProfiles(db, targetRecords)
	if err != nil {
		log.Fatalf("Failed to generate bulk data: %v", err)
	}

	duration := time.Since(start)
	recordsPerSecond := float64(targetRecords) / duration.Seconds()

	fmt.Printf("✅ Successfully inserted %d user profiles in %v\n", targetRecords, duration)
	fmt.Printf("📊 Performance: %.0f records/second\n", recordsPerSecond)
}

func setupDatabase() (*gorm.DB, error) {
	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Create PostgreSQL client using the correct function and structure
	pgClient, err := postgres.NewClient(postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
		TimeZone: cfg.Postgres.TimeZone,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Get the GORM DB instance
	db := pgClient.GormDB

	// Optimize database connection for bulk operations
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Optimize connection pool for bulk operations
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func generateBulkUserProfiles(db *gorm.DB, totalRecords int) error {
	batches := (totalRecords + BATCH_SIZE - 1) / BATCH_SIZE // Ceiling division

	for i := 0; i < batches; i++ {
		batchSize := BATCH_SIZE
		if i == batches-1 {
			// Last batch might be smaller
			batchSize = totalRecords - (i * BATCH_SIZE)
		}

		profiles := generateUserProfileBatch(batchSize, i)

		// Insert batch using GORM
		if err := db.CreateInBatches(profiles, BATCH_SIZE).Error; err != nil {
			return fmt.Errorf("failed to insert batch %d: %w", i+1, err)
		}

		// Progress reporting
		completed := (i + 1) * BATCH_SIZE
		if completed > totalRecords {
			completed = totalRecords
		}
		progress := float64(completed) / float64(totalRecords) * 100
		fmt.Printf("📈 Progress: %d/%d (%.1f%%) - Batch %d completed\n",
			completed, totalRecords, progress, i+1)
	}

	return nil
}

func generateUserProfileBatch(batchSize int, batchNumber int) []entity.UserProfile {
	profiles := make([]entity.UserProfile, batchSize)

	for i := 0; i < batchSize; i++ {
		profiles[i] = generateRandomUserProfile(batchNumber, i)
	}

	return profiles
}

func generateRandomUserProfile(batchNumber, recordIndex int) entity.UserProfile {
	// Initialize random seed with current time and unique identifiers
	r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(batchNumber*1000+recordIndex)))

	// Determine if this is a bride or groom
	isBride := r.Float32() < 0.5

	// Generate random date of birth (age between 18-50)
	minAge, maxAge := 18, 50
	ageYears := r.Intn(maxAge-minAge+1) + minAge
	dob := time.Now().AddDate(-ageYears, -r.Intn(12), -r.Intn(365))

	// Generate random height (150-185 cm)
	height := uint16(r.Intn(36) + 150) // 150-185 cm

	// Generate UNIQUE phone number using batch and record index
	phone := fmt.Sprintf("91%03d%07d", batchNumber, recordIndex*10000+r.Intn(10000))

	// Generate gender-appropriate Muslim name
	fullName := generateMuslimName(r, isBride)

	// Generate UNIQUE email using UUID for absolute uniqueness
	userUUID := uuid.New()
	emailPrefix := strings.ToLower(strings.ReplaceAll(fullName, " ", "."))
	email := fmt.Sprintf("%s.%s@testdata.com", emailPrefix, userUUID.String()[:8])

	// Random profile picture URL (gender appropriate)
	gender := "men"
	if isBride {
		gender = "women"
	}
	profilePicURL := fmt.Sprintf("https://randomuser.me/api/portraits/%s/%d.jpg", gender, r.Intn(99))

	// Generate random enum values
	community := getRandomCommunity(r)
	maritalStatus := getRandomMaritalStatus(r)
	profession := getRandomProfession(r)
	professionType := getRandomProfessionType(r)
	educationLevel := getRandomEducationLevel(r)
	homeDistrict := getRandomHomeDistrict(r)

	return entity.UserProfile{
		UserID:                userUUID,
		IsBride:               isBride,
		FullName:              fullName,
		Email:                 email,
		Phone:                 phone,
		DateOfBirth:           dob,
		HeightCm:              height,
		PhysicallyChallenged:  r.Float32() < 0.05, // 5% chance
		Community:             community,
		MaritalStatus:         maritalStatus,
		Profession:            profession,
		ProfessionType:        professionType,
		HighestEducationLevel: educationLevel,
		HomeDistrict:          homeDistrict,
		ProfileImageKey:       profilePicURL,
		LastLogin:             time.Now().Add(-time.Duration(r.Intn(30*24)) * time.Hour),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

// Random data generators using validation package enums

func getRandomCommunity(r *rand.Rand) validation.Community {
	communities := []validation.Community{
		validation.Sunni,
		validation.Mujahid,
		validation.Tabligh,
		validation.JamateIslami,
		validation.Shia,
		validation.Muslim,
	}
	return communities[r.Intn(len(communities))]
}

func getRandomMaritalStatus(r *rand.Rand) validation.MaritalStatus {
	statuses := []validation.MaritalStatus{
		validation.NeverMarried,
		validation.Divorced,
		validation.NikkahDivorce,
		validation.Widowed,
	}
	return statuses[r.Intn(len(statuses))]
}

func getRandomProfession(r *rand.Rand) validation.Profession {
	professions := []validation.Profession{
		validation.Student,
		validation.Doctor,
		validation.Engineer,
		validation.Farmer,
		validation.Teacher,
	}
	return professions[r.Intn(len(professions))]
}

func getRandomProfessionType(r *rand.Rand) validation.ProfessionType {
	types := []validation.ProfessionType{
		validation.ProfessionTypeFullTime,
		validation.ProfessionTypePartTime,
		validation.ProfessionTypeFreelance,
		validation.ProfessionTypeSelfEmployed,
		validation.ProfessionTypeNotWorking,
	}
	return types[r.Intn(len(types))]
}

func getRandomEducationLevel(r *rand.Rand) validation.EducationLevel {
	levels := []validation.EducationLevel{
		validation.LessThanHighSchool,
		validation.HighSchool,
		validation.HigherSecondary,
		validation.UnderGraduation,
		validation.PostGraduation,
	}
	return levels[r.Intn(len(levels))]
}

func getRandomHomeDistrict(r *rand.Rand) validation.HomeDistrict {
	districts := []validation.HomeDistrict{
		validation.Thiruvananthapuram,
		validation.Kollam,
		validation.Pathanamthitta,
		validation.Alappuzha,
		validation.Kottayam,
		validation.Ernakulam,
		validation.Thrissur,
		validation.Palakkad,
		validation.Malappuram,
		validation.Kozhikode,
		validation.Wayanad,
		validation.Kannur,
		validation.Kasaragod,
		validation.Idukki,
	}
	return districts[r.Intn(len(districts))]
}

// Generate gender-appropriate Muslim names
func generateMuslimName(r *rand.Rand, isBride bool) string {
	if isBride {
		return generateMuslimFemaleName(r)
	}
	return generateMuslimMaleName(r)
}

func generateMuslimMaleName(r *rand.Rand) string {
	firstNames := []string{
		"Mohammed", "Abdul", "Ahmed", "Ali", "Hassan", "Hussein", "Ibrahim", "Ismail",
		"Khalid", "Omar", "Rashid", "Tariq", "Yusuf", "Zayn", "Amjad", "Bilal",
		"Abdullah", "Abdul Rahman", "Hamza", "Umar", "Salman", "Anas", "Zaid", "Faisal",
		"Imran", "Mustafa", "Hasan", "Hussain", "Zakaria", "Sufyan", "Marwan", "Waleed",
		"Nadir", "Saeed", "Majid", "Fahad", "Adnan", "Salam", "Jalal", "Nazar",
		"Yasir", "Waseem", "Kareem", "Hakeem", "Rahim", "Rauf", "Latif", "Shafiq",
	}

	lastNames := []string{
		"Khan", "Ahmed", "Ali", "Rahman", "Hassan", "Hussein", "Ibrahim", "Malik",
		"Sheikh", "Syed", "Qureshi", "Ansari", "Hashmi", "Abbasi", "Rizvi", "Farooqui",
		"Nair", "Menon", "Kutty", "Moideen", "Haji", "Koya", "Rawther", "Thangal",
		"Chowdhury", "Hasan", "Hossain", "Islam", "Uddin", "Alam", "Sharif", "Karim",
	}

	return fmt.Sprintf("%s %s",
		firstNames[r.Intn(len(firstNames))],
		lastNames[r.Intn(len(lastNames))])
}

func generateMuslimFemaleName(r *rand.Rand) string {
	firstNames := []string{
		"Fatima", "Aisha", "Khadija", "Mariam", "Zainab", "Ruqayya", "Safiya", "Hafsa",
		"Asma", "Farah", "Layla", "Nadia", "Rania", "Sara", "Aminah", "Hafiza",
		"Samira", "Yasmin", "Rahma", "Sana", "Hiba", "Dua", "Amina", "Zahra",
		"Halima", "Sumayya", "Khawla", "Lubna", "Salma", "Hanan", "Iman", "Najla",
		"Reem", "Lina", "Maya", "Nour", "Rana", "Dina", "Hala", "Laith",
		"Sahar", "Widad", "Zara", "Hind", "Laila", "Malak", "Nada", "Razan",
	}

	lastNames := []string{
		"Khan", "Ahmed", "Ali", "Rahman", "Hassan", "Hussein", "Ibrahim", "Malik",
		"Sheikh", "Syed", "Qureshi", "Ansari", "Hashmi", "Abbasi", "Rizvi", "Farooqui",
		"Nair", "Menon", "Kutty", "Moideen", "Haji", "Koya", "Rawther", "Thangal",
		"Chowdhury", "Hasan", "Hossain", "Islam", "Uddin", "Alam", "Begum", "Khatun",
	}

	return fmt.Sprintf("%s %s",
		firstNames[r.Intn(len(firstNames))],
		lastNames[r.Intn(len(lastNames))])
}
