/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package dto

type Output struct {
	Name    string `json:"name"`
	TailKey string `json:"tail"`
	Files   []File `json:"files"`
}
type File struct {
	File string `json:"file,omitempty"`
	Size string `json:"size,omitempty"`
	Mod  string `json:"mod,omitempty"`
	Key  string `json:"key,omitempty"`
}
