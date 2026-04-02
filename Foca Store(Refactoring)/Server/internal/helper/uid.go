package helper

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"
	"voca-store/internal/domain/models"

	"gorm.io/gorm"
)

func GenerateAddressUID(tx *gorm.DB) (string, error) {
    today := time.Now().Format("20060102")
    prefix := fmt.Sprintf("ADDR-%s-", today)

    var lastAddress models.Address
    
    err := tx.Set("gorm:query_option", "FOR UPDATE").
        Unscoped().
        Where("uid LIKE ?", prefix+"%").
        Order("uid DESC").
        First(&lastAddress).Error

    nextSeq := 1

    if err == nil {
        uidStr := lastAddress.UID
        if len(uidStr) >= 4 {
            lastSeqStr := uidStr[len(uidStr)-4:]
            seq, parseErr := strconv.Atoi(lastSeqStr)
            if parseErr == nil {
                nextSeq = seq + 1
            }
        }
    } else if err != gorm.ErrRecordNotFound {
        return "", fmt.Errorf("failed to query last address UID: %w", err)
    }

    uid := fmt.Sprintf("ADDR-%s-%04d", today, nextSeq)
    return uid, nil
}

// ADDRESS UID
func NewGenerateAddressUID(tx *gorm.DB) (string, error) {
	for i := 0; i < 3; i++ {
        now := time.Now()
        random := rand.IntN(100000)
        uid := fmt.Sprintf(
            "ADDR-%s-%05d",
            now.Format("20060102-150405"),
            random,
        )
        var exist bool

        err := tx.
            Table("addresses").
            Select("count(*) > 0").
            Where("uid = ?", uid).
            Find(&exist).Error

        if err != nil {
            return "", err
        }

        if !exist {
            return uid, nil
        }
    }

    return "", fmt.Errorf("failed generate uniqe uid")
}

func GenerateCheckoutUID(tx *gorm.DB) (string, error) {
    for i := 0; i < 3; i++ {
        now := time.Now()
        random := rand.IntN(100000)
        uid := fmt.Sprintf(
            "VOCA-%s-%05d",
            now.Format("20060102-150405"),
            random,
        )
        var exists bool
        err := tx.
            Table("checkouts").
            Select("count(*) > 0").
            Where("uid = ?", uid).
            Find(&exists).Error
        if err != nil {
            return "", err
        }

        if !exists {
            return uid, nil
        }
    }

    return "", fmt.Errorf("failed generate unique uid")
}
