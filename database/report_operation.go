package database

import (
	"fmt"
	"log"
	"Forum/post_logic"
	
)

func FetchAllReports() ([]Report, error) {
	rows, err := DB.Query(`
        SELECT r.id, r.moderator_uuid, r.moderator_name, r.content, r.response, r.resolved,
               r.post_id, p.title AS post_title, p.content AS post_content
        FROM reports r
        JOIN posts p ON r.post_id = p.id
    `)
	if err != nil {
		log.Printf("Error fetching reports: %v", err) // Log the error for debugging
		return nil, err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var r Report
		if err := rows.Scan(&r.ID, &r.ModeratorUUID, &r.ModeratorName, &r.Content, &r.Response, &r.Resolved, &r.PostID, &r.PostTitle, &r.PostContent); err != nil {
			log.Printf("Error scanning report row: %v", err) // Log the error for debugging
			return nil, err
		}
		reports = append(reports, r)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over reports rows: %v", err) // Log the error for debugging
		return nil, err
	}

	return reports, nil
}

func SaveReport(report Report) error {
	query := "INSERT INTO reports (moderator_uuid, moderator_name, content, post_id, post_title, post_content) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := DB.Exec(query, report.ModeratorUUID, report.ModeratorName, report.Content, report.PostID, report.PostTitle, report.PostContent)
	if err != nil {
		return fmt.Errorf("failed to save report: %v", err)
	}
	return nil
}

func UpdateReportResponse(reportID int, response string) error {
	query := "UPDATE reports SET response = ?, resolved = TRUE WHERE id = ?"
	_, err := DB.Exec(query, response, reportID)
	if err != nil {
		return fmt.Errorf("failed to update report response: %v", err)
	}
	return nil
}

func FetchReportsByID(uuid string) (Report, error) {
	var r Report
	err := DB.QueryRow("SELECT uuid, moderator_uuid, content, response FROM reports WHERE id = ?", uuid).Scan(&r.ID, &r.ModeratorUUID, &r.Content, &r.Response)
	return r, err
}

func GetPostsByReportIDs(ids []int) ([]post_logic.Post, error) {
	var posts []post_logic.Post
	for _, id := range ids {
		p, err := FetchPostByID(id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func DeleteReportByID(reportID int) error {
	query := "DELETE FROM reports WHERE id = ?"
	_, err := DB.Exec(query, reportID)
	if err != nil {
		return fmt.Errorf("failed to delete report: %v", err)
	}
	return nil
}