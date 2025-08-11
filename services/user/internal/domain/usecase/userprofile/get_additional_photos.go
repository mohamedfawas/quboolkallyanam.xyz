package userprofile

import (
    "context"
    "fmt"

    "github.com/google/uuid"
)

func (u *userProfileUsecase) GetAdditionalPhotos(
    ctx context.Context,
    userID uuid.UUID,
) ([]string, error) {

    images, err := u.userImageRepository.ListUserImages(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to list user images: %w", err)
    }

    urls := make([]string, 0, len(images))
    for _, img := range images {
        url, err := u.photoStorage.GetDownloadURL(ctx, img.ObjectKey, u.config.MediaStorage.URLExpiry)
        if err != nil {
            return nil, fmt.Errorf("failed to generate download URL: %w", err)
        }
        urls = append(urls, url)
    }
    return urls, nil
}