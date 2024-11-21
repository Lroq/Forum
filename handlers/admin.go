package handlers

import (
	"Forum/database"
	"errors"
)

type Report struct {
	ID            int
	ModeratorID   int
	ModeratorName string
	Content       string
	Resolved      bool
	Response      string
	PostID        int
	PostTitle     string
	PostContent   string
}

func PromoteUserToModerator(userUUID string) error {
	userRole, err := database.GetUserRoleByUUID(userUUID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}
	if userRole == "moderator" {
		return errors.New("user is already a moderator")
	}
	err = database.UpdateUserRole(userUUID, "moderator")
	if err != nil {
		return errors.New("failed to promote user to moderator: " + err.Error())
	}
	return nil
}

func DemoteModeratorToUser(userUUID string) error {
	userRole, err := database.GetUserRoleByUUID(userUUID)
	if err != nil {
		return errors.New("failed to get user role: " + err.Error())
	}
	if userRole != "moderator" {
		return errors.New("user is not a moderator")
	}
	err = database.UpdateUserRole(userUUID, "utilisateur")
	if err != nil {
		return errors.New("failed to demote moderator to user: " + err.Error())
	}
	return nil
}

func ReceiveReport(report database.Report) error {
	err := database.SaveReport(report)
	if err != nil {
		return errors.New("failed to save report: " + err.Error())
	}
	return nil
}

func RespondToReport(reportID int, response string) error {
	err := database.UpdateReportResponse(reportID, response)
	if err != nil {
		return errors.New("failed to respond to report: " + err.Error())
	}
	return nil
}

func DeleteComment(commentID int) error {
	err := database.DeleteCommentByID(commentID)
	if err != nil {
		return errors.New("failed to delete comment: " + err.Error())
	}
	return nil
}
