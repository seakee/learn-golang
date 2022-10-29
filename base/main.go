package main

import "fmt"

func main() {
	// slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	// s1 := slice[2:5]
	// fmt.Println(&s1[2])
	// s2 := s1[2:6:7]
	// fmt.Println(&s2[0])
	//
	// s2 = append(s2, 100)
	// s2 = append(s2, 200)
	// fmt.Println(&s2[0])
	//
	// s1[2] = 20
	// fmt.Println(s1)
	//
	// fmt.Println(s1)
	// fmt.Println(s2)
	// fmt.Println(slice)

	// var aa []int
	// fmt.Println(len(aa))
	// fmt.Println(cap(aa))
	//
	// bb := append(aa, 1)
	//
	// fmt.Println(len(aa))
	// fmt.Println(cap(aa))
	//
	// fmt.Println(len(bb))
	// fmt.Println(cap(bb))
	//
	// cc := append(bb, 1)
	//
	// fmt.Println(len(aa))
	// fmt.Println(cap(aa))
	//
	// fmt.Println(len(bb))
	// fmt.Println(cap(bb))
	//
	// fmt.Println(len(cc))
	// fmt.Println(cap(cc))
	//
	// dd := append(cc, 1)
	//
	// fmt.Println(len(aa))
	// fmt.Println(cap(aa))
	//
	// fmt.Println(len(bb))
	// fmt.Println(cap(bb))
	//
	// fmt.Println(len(cc))
	// fmt.Println(cap(cc))
	//
	// fmt.Println(len(dd))
	// fmt.Println(cap(dd))

	// s := []int{5}
	// s = append(s, 7)
	// s = append(s, 9)
	// x := append(s, 11)
	// y := append(s, 12)
	// fmt.Println(s, x, y)

	// s := []int{1, 2}
	// s = append(s, 4, 5, 6)
	// fmt.Printf("len=%d, cap=%d", len(s), cap(s))

	// s := make([]int, 0)
	//
	// oldCap := cap(s)
	//
	// for i := 0; i < 2048; i++ {
	// 	s = append(s, i)
	//
	// 	newCap := cap(s)
	//
	// 	if newCap != oldCap {
	// 		fmt.Printf("[%d -> %4d] cap = %-4d  |  after append %-4d  cap = %-4d\n", 0, i-1, oldCap, i, newCap)
	// 		oldCap = newCap
	// 	}
	// }

	// fmt.Println(4 << (^uintptr(0) >> 63))

	a := []int{1, 2, 3}
	slice(a)
	fmt.Println("1", a)

	slicePtr1(&a)
	fmt.Println("2", a)

	slicePtr2(&a)
	fmt.Println("3", a)

	slicePtr3(&a)
	fmt.Println("3", a)
}

func slice(s []int) {
	s[0] = 10
	s = append(s, 10)
	s[1] = 10
}

func slicePtr1(s *[]int) {
	(*s)[0] = 20
	*s = append(*s, 20)
	(*s)[1] = 20
}

func slicePtr2(s *[]int) {
	b := *s
	b[0] = 30
	b = append(b, 30)
	b[1] = 30
}

func slicePtr3(s *[]int) {
	b := *s
	b = append(b, 40)
	fmt.Println("4", b)
	*s = append(*s, 50)
	fmt.Println("5", b)
}
