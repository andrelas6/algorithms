package arrays_lists

/*
 * DESCRIPTION
 *
 * Rotate Matrix: Given an image represented by an NxN matrix,
 * where each pixel in the image is represented by an integer,
 * write a method to rotate the image by 90 degrees. Can you do this in place?
 */

// Example
// [
// 	[1, 2, 3]  [7, 4, 1]
// 	[4, 5, 6]  [8, 5, 2]
//  [7, 8, 9]  [9, 6, 3]
// ]
// matrix
// a b c d
// a b c d
// a b c d
// a b c d
//
// rotated matrix 90 degrees
// a a a a
// b b b b
// c c c c
// d d d d
//
// row -> column
// row -> column
// row -> column
//
//
// Brute force
// 1 2 3: (0,0) (0,1) (0,2) -> (0,2) (1,2) (2,2)
// 4 5 6: (1,0) (1,1) (1,2) -> (0,1) (1,1) (2,1)
// 7 8 9: (2,0) (2,1) (2,2) -> (0,0) (1,0) (2,0)
//
//
//

func rotateMatrix(matrix [][]int) [][]int {

	return matrix
}
