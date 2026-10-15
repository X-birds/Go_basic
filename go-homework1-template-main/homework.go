package homework01

import (
	"sort"
)

// 1. 只出现一次的数字
// 给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
func SingleNumber(nums []int) int {
	// TODO: implement
	m := make(map[int]int)
	for _, v := range nums {
		m[v]++
	}
	for k, value := range m {
		if value == 1 {
			return k
		}
	}
	return 0
}

// 2. 回文数
// 判断一个整数是否是回文数
func IsPalindrome(x int) bool {
	// TODO: implement
	if x < 0 || (x != 0 && x%10 == 0) {
		return false
	}
	if x == 0 {
		return true
	}
	y := x
	temp := 0
	for y > 0 {
		temp = temp*10 + y%10
		y = y / 10
	}
	if x == temp {
		return true
	}
	return false
}

// 3. 有效的括号
// 给定一个只包括 '(', ')', '{', '}', '[', ']' 的字符串，判断字符串是否有效
func IsValid(s string) bool {
	// TODO: implement
	m := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}
	var stack []rune

	for _, char := range s {
		if left, ok := m[char]; ok {
			if len(stack) == 0 {
				return false
			} else if stack[len(stack)-1] != left {
				return false
			} else {
				stack = stack[:len(stack)-1]
			}
		} else {
			stack = append(stack, char)
		}
	}
	return len(stack) == 0

}

// 4. 最长公共前缀
// 查找字符串数组中的最长公共前缀
func LongestCommonPrefix(strs []string) string {
	// TODO: implement
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for i := 1; i < len(strs); i++ {
		if prefix == "" {
			return ""
		}
		c := len(prefix)
		if len(strs[i]) < c {
			c = len(strs[i])
		}
		temp := ""
		for j := 0; j < c; j++ {
			if prefix[j] == strs[i][j] {
				temp = temp + string(prefix[j])
			} else {
				break
			}
		}
		prefix = temp
	}

	if prefix != "" {
		return prefix
	}
	return ""
}

// 5. 加一
// 给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一
func PlusOne(digits []int) []int {
	l := len(digits)
	j := 1
	for i := l - 1; i >= 0 && j > 0; i-- {
		if digits[i]+j == 10 {
			digits[i] = 0
			j = 1
		} else {
			digits[i] += j
			j = 0
		}
	}
	if j == 1 {
		temp := make([]int, l+1)
		temp[0] = 1
		return temp
	}
	return digits
}

// 6. 删除有序数组中的重复项
// 给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
// 不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。
func RemoveDuplicates(nums []int) int {
	// TODO: implement
	if len(nums) == 0 {
		return 0
	}
	slow := 0
	for fast := 1; fast < len(nums); fast++ {
		if nums[slow] != nums[fast] {
			slow++
			nums[slow] = nums[fast]
		}
	}
	return slow + 1
}

// 7. 合并区间
// 以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
// 请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
func Merge(intervals [][]int) [][]int {
	// TODO: implement
	if len(intervals) == 0 {
		return nil
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	res := [][]int{}
	res = append(res, intervals[0])
	for i := 1; i < len(intervals); i++ {
		last := res[len(res)-1]
		curr := intervals[i]

		if last[1] >= curr[0] {
			if last[1] >= intervals[i][1] {
				break
			}
			res[len(res)-1][1] = intervals[i][1]
		} else {
			res = append(res, intervals[i])
		}
	}
	return res
}

// 8. 两数之和
// 给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
func TwoSum(nums []int, target int) []int {
	// TODO: implement
	if len(nums) == 0 {
		return nil
	}
	m := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		if target < nums[i] {
			break
		}
		need := target - nums[i]
		if ids, ok := m[need]; ok {
			return []int{ids, i}
		}
		m[nums[i]] = i

	}
	return nil
}
