package domain

import (
	"time"
)

// BranchNameJSON ใช้สำหรับเก็บข้อมูลชื่อสาขาในรูปแบบ JSON
type BranchNameJSON struct {
	EN string `json:"en"`
	TH string `json:"th"`
}

// BranchLocation เก็บข้อมูลที่ตั้งของสาขา
type BranchLocation struct {
	ProvinceID int `json:"province_id"`
	// สังเกต: ใน DB dump คอลัมน์ name ในตาราง branch_location เป็น float ซึ่งอาจผิดพลาด
	// หากตั้งใจให้เป็นชื่อสถานที่ ควรเปลี่ยนเป็น string
}

// Branch คือ struct หลักสำหรับข้อมูลสาขา
type Branch struct {
	ID              int64           `json:"id"`
	Name            BranchNameJSON  `json:"name"`
	Location        *BranchLocation `json:"location,omitempty"` // ใช้ pointer เพื่อให้เป็น optional
	ProductIDs      []int           `json:"product_ids,omitempty"`
	InterestIDs     []int           `json:"interest_ids,omitempty"`
	MinNormalPrice  *float64        `json:"min_normal_price,omitempty"`  // ใช้ pointer เพื่อรองรับค่า null
	MaxNormalPrice  *float64        `json:"max_normal_price,omitempty"`  // ใช้ pointer เพื่อรองรับค่า null
	MinTagthaiPrice *float64        `json:"min_tagthai_price,omitempty"` // ใช้ pointer เพื่อรองรับค่า null
	MaxTagthaiPrice *float64        `json:"max_tagthai_price,omitempty"` // ใช้ pointer เพื่อรองรับค่า null
	UpdatedAt       *time.Time      `json:"updated_at,omitempty"`        // ใช้ pointer เพื่อให้เป็น optional
}
