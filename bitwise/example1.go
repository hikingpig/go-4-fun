package main

import (
	"fmt"
)

/* link: https://medium.com/learning-the-go-programming-language/bit-hacking-with-go-e0acee258827 */

/* operators lists:
&   bitwise AND
|   bitwise OR
^   bitwise XOR
&^   AND NOT
<<   left shift
>>   right shift
*/

/*
(the input the the operator are called operands!)
Given operands a, b
AND(a, b) = 1; only if a = b = 1
               else = 0
*/
func bitwiseANDexample() {
	var x uint8 = 0xAC
	fmt.Printf("\nx = %b or %d\n", x, x)
	var y uint8 = 0xF0
	fmt.Printf("y = %b or %d\n", y, y)
	fmt.Println("----------------------x&y")
	var z uint8 = x & 0xF0
	fmt.Printf("z = %b or %d\n", z, z)
}

/* a number is even if its last bit is 0, odd if it is 1
- this is an example of AND bit masking technique, using 0000001 mask
- X AND 1 = X, X AND 0 = 0. The AND bitmasking hide the values with 0
- Any bit masked with 0 is 0, that is called "turned off bit"
*/
func checkEvenNumberExample() {
	fmt.Println()
	for x := 0; x < 10; x++ {
		if x&1 == 1 {
			fmt.Printf("%d is odd\n", x)
		} else {
			fmt.Printf("%d is even\n", x)
		}
	}
}

/*
we can choose a mask x to expose only a bit in the sequence y using AND.
- if y & x == 0, that bit is off
- if y & x > 0, that bit is on
notice that, it is on by > 0 (!=0), not == 1!
*/

func queryBitStatus() {
	var x uint8 = 4
	fmt.Printf("\nx = %b or %d\n\n", x, x)
	var y uint8 = 12
	fmt.Printf("y = %b or %d\n", y, y)
	fmt.Println("----------------------y&x")
	var z uint8 = y & x
	fmt.Printf("z = %b or %d\n", z, z)
	fmt.Println("======= z == 0?", z == 0)
	var t uint8 = 24
	fmt.Printf("\nt = %b or %d\n", t, t)
	fmt.Println("----------------------t&x")
	z = t & x
	fmt.Println("======= z == 0?", z == 0)
}

/*
Given operands a, b
OR(a, b) = 1; when a = 1 or b = 1
							else = 0
*/
func bitwiseORExample() {
	var a uint8 = 0
	fmt.Printf("\na = %b or %d\n", a, a)
	var b uint8 = 196
	fmt.Printf("b = %b or %d\n", b, b)
	fmt.Println("--------------------- a |=196")
	a |= 196
	fmt.Printf("a = %b or %d\n", a, a)
}

/*Using biwise OR principle: Y OR 1 = 1 and Y OR 0 = Y. OR will mask the value we want to hide with 1
- Any bit masked with 1 will be 1, that is called "turned on bit"
*/

func turnOnBitExample() {
	var a uint8 = 195
	fmt.Printf("\na = %b or %d\n", a, a)
	var b uint8 = 3
	fmt.Printf("b = %b or %d\n", b, b)
	fmt.Println("----------------- a |=b")
	a |= b
	fmt.Printf("a = %b or %d\n", a, a)
}

/* we can combine both AND and OR for multiple selection
- the selections are coded with 2^n integers
- the multiple selection using OR m = a | b | c (we choose both a, b, c)
- the check using AND bit masking to check what we chose:
	- m & a == 0 -> are is not chosen. m & a != 0, a chosen
	- similar to b, c, d ..

- byte is equivalent to uint8!
*/
func multipleSelection() {
	const MARRY = 2
	const JENNIFER = 4
	const YVONNE = 8
	const JACK = 16
	var whoilove = func(lovers byte) {
		if lovers&MARRY != 0 {
			fmt.Println("I love Mary more than anything")
		}
		if lovers&JENNIFER != 0 {
			fmt.Println("I love Jennifer more than anything")
		}
		if lovers&YVONNE != 0 {
			fmt.Println("I love Yvonne more than anything")
		}
		if lovers&JACK != 0 {
			fmt.Println("It's a mistake, check the input")
		}
	}
	whoilove(MARRY | JENNIFER | YVONNE)
}

/*
Given operands a, b
XOR(a, b) = 1; only if a != b
     else = 0
*/

func bitwiseXORExample() {
	var a uint16 = 0xCEFF
	var b uint16 = 0xFF00
	fmt.Printf("\nx = %b or %d\n", a, a)
	fmt.Printf("y = %b or %d\n", b, b)
	fmt.Println("----------------------x^y")
	var z uint16 = a ^ b
	// use %016b format to print all the prefixed 0
	fmt.Printf("z = %016b or %d\n", z, z)
}

/*
Y XOR 1 will be 1 if Y = 0, 0 if Y = 1.
The bit Y is toggled with Y XOR 1
*/
func toggleBits() {
	var a uint8 = 14
	// use %08b format to print out the prefixed 0
	fmt.Printf("\nx = %08b or %d\n", a, a)
	var b uint8 = 30
	fmt.Printf("y = %08b or %d\n", b, b)
	var z uint8 = a ^ b
	fmt.Printf("z = %08b or %d\n", z, z)
}

/* the two's complement number of a is ^a, short for 11111111 ^a. 1 ^ a just not reverse all the bits!
To get the negative number of a:
-a = ^a + 1
*/
func negativeNumber() {
	// replace uint8 or byte with int8 to have a signed integer
	var a int8 = 0x0F
	fmt.Printf("a = %08b or %d\n", a, a)
	fmt.Printf("^a = %08b or %d\n", ^a, ^a)
	b := ^a + 1
	fmt.Printf("-a = %08b or %d\n", b, b)
	fmt.Printf("-b = %08b or %d\n", ^b+1, ^b+1)
}

/* the first bit in signed integer representation also defines sign
	1 -> negative
	0 -> positive
so if a negative number XOR a positive number, this first bit is always 1. that means we always have a negative number. the rest bits in the sequence are not important.
*/
func checkSameSign() {
	a, b := -12, 25
	fmt.Printf("\nx = %08b or %d\n", a, a)
	fmt.Printf("y = %08b or %d\n", b, b)
	c := ^a
	fmt.Printf("c = %08b or %d\n", c, c)
	d := ^c
	fmt.Printf("d = %08b or %d\n", d, d)
	fmt.Println("a and b have same sign?", (a^b) >= 0)
}

/*
Given operands a, b
AND_NOT(a, b) = AND(a, NOT(b))

Here we will clear (hide with 0) some bit defined by a mask "b"
- The bits in mask value 1 at positions to be cleared
- We then use XOR to reverse b's bits to an AND mask where the bits have values 0 in position to be hidden
- Apply AND to get the masked value.
*/

func bitwiseANDNOTExample() {
	var a byte = 0xAB
	fmt.Printf("a = %08b or %d\n", a, a)
	var b byte = 0x0F
	fmt.Printf("b = %08b or %d\n", b, b)
	fmt.Printf("^b = %08b or %d\n", ^b, ^b)

	a &^= b
	fmt.Println("--------------------------- a &^ b")
	fmt.Printf("a = %08b or %d\n", a, a)
}

/*
Given integer operands a and n,
a << n; shifts all bits in a to the left n times
a >> n; shifts all bits in a to the right n times

right shift << means multiply 2^n times
left shift >> means dividing 2^n times (only valid if has n 0 bits at the end)
*/
func bitwiseShiftExample() {
	var a int8 = 3
	fmt.Printf("a = %08b or %d\n", a, a)
	fmt.Printf("a <<1 = %08b or %d\n", a<<1, a<<1)
	fmt.Printf("a <<2 = %08b or %d\n", a<<2, a<<2)
	fmt.Printf("a <<3 = %08b or %d\n", a<<3, a<<3)

	a = 120
	fmt.Printf("\na = %08b or %d\n", a, a)
	fmt.Printf("a>>1 = %08b or %d\n", a>>1, a>>1)
	fmt.Printf("a>>2 = %08b or %d\n", a>>2, a>>2)
}

/*
we can shift 1 to a position and use that bit to manipulate or test a bit in another sequence
- manipulate: |: turn on bit, &^: unset
- &	: test
*/
func manipulateBit() {
	/* turn on a bit */
	var a int8 = 8
	fmt.Printf("a = %08b or %d\n", a, a)
	fmt.Printf("1<<2 = %008b or %d\n", 1<<2, 1<<2)
	a = a | (1 << 2)
	fmt.Println("-------------------- a | (1 << 2)")
	fmt.Printf("a = %08b or %d\n", a, a)

	/* test a bit */
	a = 12
	if a&(1<<2) != 0 {
		fmt.Println("\ntake action")
	}

	/* unset (turn off a bit) */
	a = 13
	fmt.Printf("\n%04b\n", a)
	a = a &^ (1 << 2)
	fmt.Printf("%04b\n", a)
}

func main() {
	// bitwiseANDexample()
	// queryBitStatus()
	// checkEvenNumberExample()
	// multipleSelection()
	// bitwiseORExample()
	// turnOnBitExample()
	// bitwiseXORExample()
	// toggleBits()
	// negativeNumber()
	// checkSameSign()
	// bitwiseANDNOTExample()
	// bitwiseShiftExample()
	manipulateBit()
}

/* ip "inverse mask" is another usage mask. It creates the access control lists of ip.
- network address: 10.1.1.0
- mask:						 0.0.0.255
- allowed addresses: 10.1.1.x (x from 0 to 255)
We see clearer using binary:
- network address: 00001010.00000001.00000001.00000000
- mask:						 00000000.00000000.00000000.11111111
- allowed addresses: 00001010.00000001.00000001.x

the 00000000 means the octet must be the same that as that in the network addresss. the 11111111 means the octet can be any value.

Not sure how the mask will work to check the address. See later.

the "inverse mask" is obtained by subtracting the normal mask from all 255s
							255.255.255.255
normal mask		255.255.255.0
inverse mask	0.0.0.255
*/
