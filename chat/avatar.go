package main

import (
	"errors"
	"io/ioutil"
	"path"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client, // or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get // a URL for the specified client.
	GetAvatarURL(ChatUser) (string, error)
}

// TryAvatars type is simply a slice of Avatar objects
// that we are free to add methods to
type TryAvatars []Avatar

// GetAvatarURL gets the avatar URL for the slice of avatars,
// or returns an error if something goes wrong.
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}

	return "", ErrNoAvatarURL
}

// AuthAvatar type implements Avatar interface
type AuthAvatar struct{}

// UseAuthAvatar has the AuthAvatar type
var UseAuthAvatar AuthAvatar

// GetAvatarURL gets the avatar URL for the specified client,
// or returns an error if something goes wrong.
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()

	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}

	return url, nil
}

// GravatarAvatar type implements Avatar interface
type GravatarAvatar struct{}

// UseGravatarAvatar has the AuthAvatar type
var UseGravatarAvatar GravatarAvatar

// GetAvatarURL gets the avatar URL for the specified client,
// or returns an error if something goes wrong.
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar type implements Avatar interface
type FileSystemAvatar struct{}

// UseFileSystemAvatar has the AuthAvatar type
var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL gets the avatar URL for the specified client,
// or returns an error if something goes wrong.
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}

			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}

	return "", ErrNoAvatarURL
}
