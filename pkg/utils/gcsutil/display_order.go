package gcsutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/constants"
)

func ExtractDisplayOrder(objectKey string) (int32, error) {
	parts := strings.Split(objectKey, "/")
	if len(parts) != 3 || parts[0] != constants.AdditionalPhotoStorageDirectory {
		return 0, fmt.Errorf("object key is not in additional-photo format: %q", objectKey)
	}

	order64, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid displayOrder %q in objectKey: %w", parts[2], err)
	}

	return int32(order64), nil
}