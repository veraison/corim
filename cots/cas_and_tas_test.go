// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package cots

import (
	"bytes"
	//"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	//"time"
	//
	//"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
	//"github.com/veraison/corim/comid"
)

var (
	// Good CA from PKITS
	ca = []byte{
		0x30, 0x82, 0x03, 0x7C, 0x30, 0x82, 0x02, 0x64, 0xA0, 0x03, 0x02, 0x01, 0x02, 0x02, 0x01, 0x02,
		0x30, 0x0D, 0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x01, 0x0B, 0x05, 0x00, 0x30,
		0x45, 0x31, 0x0B, 0x30, 0x09, 0x06, 0x03, 0x55, 0x04, 0x06, 0x13, 0x02, 0x55, 0x53, 0x31, 0x1F,
		0x30, 0x1D, 0x06, 0x03, 0x55, 0x04, 0x0A, 0x13, 0x16, 0x54, 0x65, 0x73, 0x74, 0x20, 0x43, 0x65,
		0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x20, 0x32, 0x30, 0x31, 0x31, 0x31,
		0x15, 0x30, 0x13, 0x06, 0x03, 0x55, 0x04, 0x03, 0x13, 0x0C, 0x54, 0x72, 0x75, 0x73, 0x74, 0x20,
		0x41, 0x6E, 0x63, 0x68, 0x6F, 0x72, 0x30, 0x1E, 0x17, 0x0D, 0x31, 0x30, 0x30, 0x31, 0x30, 0x31,
		0x30, 0x38, 0x33, 0x30, 0x30, 0x30, 0x5A, 0x17, 0x0D, 0x33, 0x30, 0x31, 0x32, 0x33, 0x31, 0x30,
		0x38, 0x33, 0x30, 0x30, 0x30, 0x5A, 0x30, 0x40, 0x31, 0x0B, 0x30, 0x09, 0x06, 0x03, 0x55, 0x04,
		0x06, 0x13, 0x02, 0x55, 0x53, 0x31, 0x1F, 0x30, 0x1D, 0x06, 0x03, 0x55, 0x04, 0x0A, 0x13, 0x16,
		0x54, 0x65, 0x73, 0x74, 0x20, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
		0x73, 0x20, 0x32, 0x30, 0x31, 0x31, 0x31, 0x10, 0x30, 0x0E, 0x06, 0x03, 0x55, 0x04, 0x03, 0x13,
		0x07, 0x47, 0x6F, 0x6F, 0x64, 0x20, 0x43, 0x41, 0x30, 0x82, 0x01, 0x22, 0x30, 0x0D, 0x06, 0x09,
		0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x01, 0x01, 0x05, 0x00, 0x03, 0x82, 0x01, 0x0F, 0x00,
		0x30, 0x82, 0x01, 0x0A, 0x02, 0x82, 0x01, 0x01, 0x00, 0x90, 0x58, 0x9A, 0x47, 0x62, 0x8D, 0xFB,
		0x5D, 0xF6, 0xFB, 0xA0, 0x94, 0x8F, 0x7B, 0xE5, 0xAF, 0x7D, 0x39, 0x73, 0x20, 0x6D, 0xB5, 0x59,
		0x0E, 0xCC, 0xC8, 0xC6, 0xC6, 0xB4, 0xAF, 0xE6, 0xF2, 0x67, 0xA3, 0x0B, 0x34, 0x7A, 0x73, 0xE7,
		0xFF, 0xA4, 0x98, 0x44, 0x1F, 0xF3, 0x9C, 0x0D, 0x23, 0x2C, 0x5E, 0xAF, 0x21, 0xE6, 0x45, 0xDA,
		0x04, 0x6A, 0x96, 0x2B, 0xEB, 0xD2, 0xC0, 0x3F, 0xCF, 0xCE, 0x9E, 0x4E, 0x60, 0x6A, 0x6D, 0x5E,
		0x61, 0x8F, 0x72, 0xD8, 0x43, 0xB4, 0x0C, 0x25, 0xAD, 0xA7, 0xE4, 0x18, 0xE4, 0xB8, 0x1A, 0xA2,
		0x09, 0xF3, 0xE9, 0x3D, 0x5C, 0x62, 0xAC, 0xFA, 0xF4, 0x14, 0x5C, 0x92, 0xAC, 0x3A, 0x4E, 0x3B,
		0x46, 0xEC, 0xC3, 0xE8, 0xF6, 0x6E, 0xA6, 0xAE, 0x2C, 0xD7, 0xAC, 0x5A, 0x2D, 0x5A, 0x98, 0x6D,
		0x40, 0xB6, 0xE9, 0x47, 0x18, 0xD3, 0xC1, 0xA9, 0x9E, 0x82, 0xCD, 0x1C, 0x96, 0x52, 0xFC, 0x49,
		0x97, 0xC3, 0x56, 0x59, 0xDD, 0xDE, 0x18, 0x66, 0x33, 0x65, 0xA4, 0x8A, 0x56, 0x14, 0xD1, 0xE7,
		0x50, 0x69, 0x9D, 0x88, 0x62, 0x97, 0x50, 0xF5, 0xFF, 0xF4, 0x7D, 0x1F, 0x56, 0x32, 0x00, 0x69,
		0x0C, 0x23, 0x9C, 0x60, 0x1B, 0xA6, 0x0C, 0x82, 0xBA, 0x65, 0xA0, 0xCC, 0x8C, 0x0F, 0xA5, 0x7F,
		0x84, 0x94, 0x53, 0x94, 0xAF, 0x7C, 0xFB, 0x06, 0x85, 0x67, 0x14, 0xA8, 0x48, 0x5F, 0x37, 0xBE,
		0x56, 0x64, 0x06, 0x49, 0x6C, 0x59, 0xC6, 0xF5, 0x83, 0x50, 0xDF, 0x74, 0x52, 0x5D, 0x2D, 0x2C,
		0x4A, 0x4B, 0x82, 0x4D, 0xCE, 0x57, 0x15, 0x01, 0xE1, 0x55, 0x06, 0xB9, 0xFD, 0x79, 0x38, 0x93,
		0xA9, 0x82, 0x8D, 0x71, 0x89, 0xB2, 0x0D, 0x3E, 0x65, 0xAD, 0xD7, 0x85, 0x5D, 0x6B, 0x63, 0x7D,
		0xCA, 0xB3, 0x4A, 0x96, 0x82, 0x46, 0x64, 0xDA, 0x8B, 0x02, 0x03, 0x01, 0x00, 0x01, 0xA3, 0x7C,
		0x30, 0x7A, 0x30, 0x1F, 0x06, 0x03, 0x55, 0x1D, 0x23, 0x04, 0x18, 0x30, 0x16, 0x80, 0x14, 0xE4,
		0x7D, 0x5F, 0xD1, 0x5C, 0x95, 0x86, 0x08, 0x2C, 0x05, 0xAE, 0xBE, 0x75, 0xB6, 0x65, 0xA7, 0xD9,
		0x5D, 0xA8, 0x66, 0x30, 0x1D, 0x06, 0x03, 0x55, 0x1D, 0x0E, 0x04, 0x16, 0x04, 0x14, 0x58, 0x01,
		0x84, 0x24, 0x1B, 0xBC, 0x2B, 0x52, 0x94, 0x4A, 0x3D, 0xA5, 0x10, 0x72, 0x14, 0x51, 0xF5, 0xAF,
		0x3A, 0xC9, 0x30, 0x0E, 0x06, 0x03, 0x55, 0x1D, 0x0F, 0x01, 0x01, 0xFF, 0x04, 0x04, 0x03, 0x02,
		0x01, 0x06, 0x30, 0x17, 0x06, 0x03, 0x55, 0x1D, 0x20, 0x04, 0x10, 0x30, 0x0E, 0x30, 0x0C, 0x06,
		0x0A, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x02, 0x01, 0x30, 0x01, 0x30, 0x0F, 0x06, 0x03, 0x55,
		0x1D, 0x13, 0x01, 0x01, 0xFF, 0x04, 0x05, 0x30, 0x03, 0x01, 0x01, 0xFF, 0x30, 0x0D, 0x06, 0x09,
		0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x01, 0x0B, 0x05, 0x00, 0x03, 0x82, 0x01, 0x01, 0x00,
		0x35, 0x87, 0x97, 0x16, 0xE6, 0x75, 0x35, 0xCD, 0xC0, 0x12, 0xFF, 0x96, 0x5C, 0x21, 0x42, 0xAC,
		0x27, 0x6B, 0x32, 0xBB, 0x08, 0x2D, 0x96, 0xB1, 0x70, 0x41, 0xAA, 0x03, 0x4F, 0x5A, 0x3E, 0xE6,
		0xB6, 0xF4, 0x3E, 0x68, 0xB1, 0xBC, 0xFF, 0x9D, 0x10, 0x73, 0x64, 0xAE, 0x9F, 0xBA, 0x36, 0x56,
		0x7C, 0x05, 0xF4, 0x3D, 0x7C, 0x51, 0x47, 0xBC, 0x1A, 0x3D, 0xEE, 0x3D, 0x46, 0x07, 0xFA, 0x84,
		0x88, 0xD6, 0xF0, 0xDD, 0xC8, 0xA7, 0x23, 0x98, 0xC6, 0xCA, 0x45, 0x4E, 0x2B, 0x93, 0x47, 0xA8,
		0xDD, 0x41, 0xCD, 0x0D, 0x7C, 0x2A, 0x21, 0x57, 0x3D, 0x09, 0x04, 0xBD, 0xB2, 0x6C, 0x95, 0xFB,
		0x1D, 0x47, 0x0B, 0x02, 0xF8, 0x4D, 0x3A, 0xEA, 0xF8, 0xB5, 0xCB, 0x2B, 0x1F, 0xEA, 0x56, 0x28,
		0xF4, 0x62, 0xA9, 0x3E, 0x50, 0x97, 0xC0, 0xB6, 0xB8, 0x36, 0x8E, 0x76, 0x0A, 0x5E, 0xC0, 0xAE,
		0x14, 0xC0, 0x50, 0x42, 0x75, 0x82, 0x1A, 0xBC, 0x1A, 0xD6, 0x0D, 0x53, 0xA6, 0x14, 0x69, 0xFD,
		0x19, 0x98, 0x1E, 0x73, 0x32, 0x9D, 0x81, 0x66, 0x66, 0xB5, 0xED, 0xCC, 0x5C, 0xFE, 0x53, 0xD5,
		0xC4, 0x03, 0xB0, 0xBE, 0x80, 0xFA, 0xB8, 0x92, 0xA0, 0xC8, 0xFE, 0x25, 0x5F, 0x21, 0x3D, 0x6C,
		0xEA, 0x50, 0x6D, 0x74, 0x1E, 0x74, 0x96, 0xB0, 0xD5, 0xC2, 0x5D, 0xA8, 0x61, 0xF0, 0x2F, 0x5B,
		0xFE, 0xAC, 0x0B, 0x6B, 0x1E, 0xD9, 0x09, 0x5E, 0x66, 0x27, 0x54, 0x9A, 0xBC, 0xE2, 0x54, 0xD3,
		0xF8, 0xA0, 0x47, 0x97, 0x20, 0xDA, 0x24, 0x53, 0xA4, 0xFA, 0xA7, 0xFF, 0xC7, 0x33, 0x51, 0x46,
		0x41, 0x8C, 0x36, 0x8C, 0xEB, 0xE9, 0x29, 0xC2, 0xAD, 0x58, 0x24, 0x80, 0x9D, 0xE8, 0x04, 0x6E,
		0x0B, 0x06, 0x63, 0x30, 0x13, 0x2A, 0x39, 0x8F, 0x24, 0xF2, 0x74, 0x9E, 0x91, 0xC5, 0xAB, 0x33,
	}

	// Trust anchor from PKITS
	ta = []byte{
		0x30, 0x82, 0x03, 0x47, 0x30, 0x82, 0x02, 0x2F, 0xA0, 0x03, 0x02, 0x01, 0x02, 0x02, 0x01, 0x01,
		0x30, 0x0D, 0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x01, 0x0B, 0x05, 0x00, 0x30,
		0x45, 0x31, 0x0B, 0x30, 0x09, 0x06, 0x03, 0x55, 0x04, 0x06, 0x13, 0x02, 0x55, 0x53, 0x31, 0x1F,
		0x30, 0x1D, 0x06, 0x03, 0x55, 0x04, 0x0A, 0x13, 0x16, 0x54, 0x65, 0x73, 0x74, 0x20, 0x43, 0x65,
		0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x20, 0x32, 0x30, 0x31, 0x31, 0x31,
		0x15, 0x30, 0x13, 0x06, 0x03, 0x55, 0x04, 0x03, 0x13, 0x0C, 0x54, 0x72, 0x75, 0x73, 0x74, 0x20,
		0x41, 0x6E, 0x63, 0x68, 0x6F, 0x72, 0x30, 0x1E, 0x17, 0x0D, 0x31, 0x30, 0x30, 0x31, 0x30, 0x31,
		0x30, 0x38, 0x33, 0x30, 0x30, 0x30, 0x5A, 0x17, 0x0D, 0x33, 0x30, 0x31, 0x32, 0x33, 0x31, 0x30,
		0x38, 0x33, 0x30, 0x30, 0x30, 0x5A, 0x30, 0x45, 0x31, 0x0B, 0x30, 0x09, 0x06, 0x03, 0x55, 0x04,
		0x06, 0x13, 0x02, 0x55, 0x53, 0x31, 0x1F, 0x30, 0x1D, 0x06, 0x03, 0x55, 0x04, 0x0A, 0x13, 0x16,
		0x54, 0x65, 0x73, 0x74, 0x20, 0x43, 0x65, 0x72, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x65,
		0x73, 0x20, 0x32, 0x30, 0x31, 0x31, 0x31, 0x15, 0x30, 0x13, 0x06, 0x03, 0x55, 0x04, 0x03, 0x13,
		0x0C, 0x54, 0x72, 0x75, 0x73, 0x74, 0x20, 0x41, 0x6E, 0x63, 0x68, 0x6F, 0x72, 0x30, 0x82, 0x01,
		0x22, 0x30, 0x0D, 0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 0x01, 0x01, 0x01, 0x05, 0x00,
		0x03, 0x82, 0x01, 0x0F, 0x00, 0x30, 0x82, 0x01, 0x0A, 0x02, 0x82, 0x01, 0x01, 0x00, 0xB9, 0x99,
		0x51, 0x89, 0x11, 0x47, 0xE0, 0xCD, 0x45, 0xB9, 0x84, 0x27, 0x82, 0x13, 0x02, 0x16, 0xCD, 0x44,
		0x39, 0xAA, 0xAC, 0xDB, 0x09, 0xC3, 0xDE, 0xE2, 0x2C, 0x4E, 0xDB, 0xA4, 0x57, 0x9F, 0x94, 0x35,
		0x28, 0x35, 0xC2, 0x6B, 0x64, 0xF6, 0x5B, 0x0B, 0x6D, 0x3C, 0x6C, 0x6F, 0xA0, 0xE4, 0xCD, 0x00,
		0x4F, 0x35, 0x74, 0xBC, 0xEA, 0xA0, 0xF3, 0xC0, 0xEC, 0xD7, 0xF1, 0x63, 0x14, 0x32, 0xB5, 0xA0,
		0xF2, 0x2D, 0xAA, 0xC8, 0x83, 0x22, 0x62, 0x48, 0xAB, 0x4C, 0x5B, 0xFD, 0xEB, 0x79, 0xC3, 0xBD,
		0x96, 0x34, 0xFC, 0x47, 0x56, 0xB7, 0x2C, 0xAF, 0xB0, 0x29, 0xE8, 0x31, 0xDF, 0x77, 0x02, 0xE9,
		0x34, 0xC9, 0xDC, 0xAA, 0xDC, 0xD7, 0xF7, 0x68, 0x54, 0xFE, 0x21, 0x95, 0x1C, 0xB1, 0x3F, 0xC3,
		0xF3, 0x82, 0x5B, 0x00, 0x08, 0x21, 0xB4, 0xC7, 0x6B, 0x17, 0x6B, 0x18, 0xC5, 0x06, 0xEC, 0x39,
		0x09, 0x27, 0xA8, 0x88, 0xEE, 0x6C, 0xEA, 0xA5, 0xCC, 0x8F, 0xBF, 0x20, 0x00, 0xA3, 0xD3, 0xB7,
		0x48, 0x78, 0xC9, 0xAF, 0xEB, 0xB0, 0xE6, 0xF4, 0xAB, 0x1D, 0x1A, 0xDE, 0xB6, 0x75, 0x76, 0xBA,
		0x7D, 0x1B, 0xA2, 0x1B, 0xC6, 0xB2, 0x53, 0x7A, 0xE0, 0xC6, 0x3F, 0x50, 0x88, 0x91, 0x9C, 0xF1,
		0x70, 0x77, 0x68, 0x03, 0xEF, 0xA6, 0xF2, 0x0F, 0x3A, 0x0C, 0xD6, 0x2A, 0x32, 0x2D, 0x10, 0xA5,
		0x95, 0xF0, 0x49, 0xE6, 0xC4, 0x43, 0xCF, 0xDC, 0x6B, 0x50, 0x73, 0x62, 0x81, 0x30, 0x14, 0x76,
		0x56, 0xE5, 0x6F, 0xAC, 0xAB, 0x9C, 0xDB, 0x4D, 0x26, 0x69, 0x2B, 0x44, 0xE6, 0x2F, 0x92, 0x1E,
		0xAC, 0x2C, 0x44, 0xD4, 0x87, 0xD1, 0x18, 0x87, 0x40, 0x2B, 0xFB, 0xDE, 0x1F, 0x65, 0x55, 0xDF,
		0x19, 0x58, 0xB9, 0xED, 0x2C, 0x3B, 0xC7, 0x58, 0xF2, 0x00, 0xD4, 0xE5, 0x03, 0x8D, 0x02, 0x03,
		0x01, 0x00, 0x01, 0xA3, 0x42, 0x30, 0x40, 0x30, 0x1D, 0x06, 0x03, 0x55, 0x1D, 0x0E, 0x04, 0x16,
		0x04, 0x14, 0xE4, 0x7D, 0x5F, 0xD1, 0x5C, 0x95, 0x86, 0x08, 0x2C, 0x05, 0xAE, 0xBE, 0x75, 0xB6,
		0x65, 0xA7, 0xD9, 0x5D, 0xA8, 0x66, 0x30, 0x0E, 0x06, 0x03, 0x55, 0x1D, 0x0F, 0x01, 0x01, 0xFF,
		0x04, 0x04, 0x03, 0x02, 0x01, 0x06, 0x30, 0x0F, 0x06, 0x03, 0x55, 0x1D, 0x13, 0x01, 0x01, 0xFF,
		0x04, 0x05, 0x30, 0x03, 0x01, 0x01, 0xFF, 0x30, 0x0D, 0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7,
		0x0D, 0x01, 0x01, 0x0B, 0x05, 0x00, 0x03, 0x82, 0x01, 0x01, 0x00, 0x98, 0xA1, 0xAF, 0x6E, 0x47,
		0x9E, 0x4A, 0x25, 0x39, 0x29, 0xC3, 0x23, 0xE2, 0x84, 0x88, 0x17, 0x1E, 0xAE, 0xFF, 0x67, 0xC7,
		0x71, 0xDE, 0xA6, 0x65, 0x08, 0x94, 0x30, 0x13, 0x9A, 0x90, 0x05, 0x95, 0xC8, 0xB1, 0xFE, 0x5E,
		0x0B, 0xC5, 0x5A, 0xFD, 0x08, 0xE7, 0x4C, 0x73, 0x82, 0xE0, 0x6B, 0x78, 0x33, 0x0A, 0x67, 0xAA,
		0xA9, 0xA8, 0x0E, 0xAC, 0xAA, 0x49, 0x6F, 0x29, 0x05, 0x54, 0x20, 0x01, 0x41, 0x40, 0x5E, 0xA3,
		0xBD, 0xD6, 0x72, 0xC2, 0x41, 0xBE, 0x3E, 0xFF, 0x27, 0xE3, 0x8A, 0x23, 0x63, 0xCE, 0x9A, 0xE9,
		0xC5, 0x09, 0x5C, 0xA8, 0x86, 0x58, 0x8C, 0x94, 0x95, 0xB5, 0xC9, 0x07, 0xD4, 0x80, 0xD6, 0x14,
		0xB1, 0x5E, 0xB8, 0x65, 0xBE, 0x3A, 0x03, 0xCF, 0x10, 0x58, 0x0D, 0xE2, 0x18, 0xD4, 0xE5, 0x81,
		0x48, 0xB0, 0x4D, 0xAB, 0xD6, 0x2E, 0x3E, 0x8B, 0x19, 0x1E, 0xBA, 0x8A, 0xF1, 0xC6, 0xC2, 0xB5,
		0x0F, 0xF8, 0xD8, 0x6D, 0x5C, 0x10, 0x31, 0x72, 0xEA, 0x9F, 0x5A, 0x63, 0x26, 0xDB, 0x4E, 0x3C,
		0x9C, 0x2E, 0x03, 0xCF, 0xA1, 0xA3, 0x57, 0xF6, 0x73, 0xC1, 0x6B, 0x2A, 0x5A, 0xA3, 0x1E, 0x03,
		0xB0, 0xC6, 0xE1, 0xE1, 0xB2, 0x21, 0x8D, 0xE8, 0xC2, 0xA0, 0xB3, 0x56, 0xDA, 0x6A, 0x5A, 0x51,
		0xFE, 0x59, 0xCD, 0x14, 0x22, 0x29, 0xAB, 0xEF, 0xFE, 0xDD, 0xC9, 0xE1, 0xB9, 0xF0, 0xE3, 0xBF,
		0x13, 0x32, 0xE6, 0x58, 0x3E, 0x73, 0x08, 0x0C, 0x0A, 0x21, 0xB3, 0x0B, 0x19, 0xF8, 0x9F, 0x87,
		0x41, 0x33, 0x0E, 0x35, 0x0B, 0xEE, 0x0A, 0x84, 0xD1, 0x7B, 0x66, 0xD2, 0xAE, 0x29, 0x26, 0x75,
		0x79, 0xCC, 0xF3, 0xB5, 0x70, 0xFD, 0x35, 0x49, 0x06, 0x50, 0x6C, 0x37, 0x2F, 0x3A, 0x4B, 0x0C,
		0x96, 0xB3, 0xCF, 0x72, 0x61, 0x95, 0x9F, 0xFD, 0xA5, 0x7C, 0xC5,
	}
)

func TestTasAndCas(t *testing.T) {
	tas := NewTrustAnchor()
	tas.Data = ta
	tas.Format = TaFormatCertificate
	ct, _ := tas.ToCBOR()
	assert.NotNil(t, ct)

	tv := NewTasAndCas()
	tv.AddTaCert(ta)

	assert.Nil(t, tv.Valid())

	c, _ := tv.ToCBOR()
	j, _ := tv.ToJSON()
	assert.NotNil(t, j)
	assert.NotNil(t, c)

	tv2 := NewTasAndCas()
	tv2.FromJSON(j)

	tv3 := NewTasAndCas()
	tv3.FromCBOR(c)

	assert.Truef(t, len(tv2.Tas[0].Data) == len(ta) && 0 == bytes.Compare(ta, tv2.Tas[0].Data), "Compare TA value")
	assert.Truef(t, TaFormatCertificate == tv2.Tas[0].Format, "Compare TA value")
}

func TestEmptyTasAndCas(t *testing.T) {
	tv := NewTasAndCas()

	assert.NotNil(t, tv.Valid())
}

func TestEmptyTas(t *testing.T) {
	tv := NewTasAndCas()
	tv.AddCaCert(ca)

	assert.NotNil(t, tv.Valid())
}

func TestEmptyCas(t *testing.T) {
	tv := NewTasAndCas()
	tv.AddTaCert(ta)

	assert.Nil(t, tv.Valid())
}
