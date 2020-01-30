package main

import (
	"errors"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client, // or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get // a URL for the specified client.
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar type implements Avatar interface
type AuthAvatar struct{}

// UseAuthAvatar has the AuthAvatar type
var UseAuthAvatar AuthAvatar

// GetAvatarURL gets the avatar URL for the specified client,
// or returns an error if something goes wrong.
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

// GravatarAvatar type implements Avatar interface
type GravatarAvatar struct{}

// UseGravatarAvatar has the AuthAvatar type
var UseGravatarAvatar GravatarAvatar

// GetAvatarURL gets the avatar URL for the specified client,
// or returns an error if something goes wrong.
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

// FileSystemAvatar type implements Avatar interface
type FileSystemAvatar struct{}

// UseFileSystemAvatar has the AuthAvatar type
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL gets the avatar URL for the specified client,
// or returns an error if something goes wrong.
func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "/avatars/" + useridStr + ".jpg", nil
		}
	}

	return "", ErrNoAvatarURL
}
