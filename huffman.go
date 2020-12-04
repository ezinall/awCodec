package main

const (
	maxTableEntry = 15 // Maximum Huffman table entry index
)

var huffmanCodes = [][][2]int{
	// Table 0 Not used
	// 1
	{{0b1, 1}, {0b001, 3}},
	{{0b01, 2}, {0b00, 3}},
	// 2
	{{0b1, 1}, {0b010, 3}, {0b000001, 6}},
	{{0b011, 3}, {0b001, 3}, {0b00001, 5}},
	{{0b00011, 5}, {0b00010, 5}, {0b000000, 6}},
	// 3
	{{0b11, 2}, {0b10, 2}, {0b000001, 6}},
	{{0b001, 3}, {0b01, 2}, {0b00001, 5}},
	{{0b00011, 5}, {0b00010, 5}, {0b000000, 6}},
	// 4 Not used
	// 5
	{{0b1, 1}, {0b010, 3}, {0b000110, 6}, {0b0000101, 7}},
	{{0b011, 3}, {0b001, 3}, {0b000100, 6}, {0b0000100, 7}},
	{{0b000111, 6}, {0b000101, 6}, {0b0000111, 7}, {0b00000001, 8}},
	{{0b0000110, 7}, {0b000001, 6}, {0b0000001, 7}, {0b00000000, 8}},
	// 6
	{{0b111, 3}, {0b011, 3}, {0b00101, 5}, {0b0000001, 7}},
	{{0b110, 3}, {0b10, 2}, {0b0011, 4}, {0b00010, 5}},
	{{0b0101, 4}, {0b0100, 4}, {0b00100, 5}, {0b000001, 6}},
	{{0b000011, 6}, {0b00011, 5}, {0b000010, 6}, {0b0000000, 7}},
	// 7
	{{0b1, 3}, {0b010, 3}, {0b001010, 6}, {0b00010011, 8}, {0b00010000, 8}, {0b000001010, 9}},
	{{0b011, 3}, {0b0011, 4}, {0b000111, 6}, {0b0001010, 7}, {0b0000101, 7}, {0b00000011, 8}},
	{{0b001011, 6}, {0b00100, 5}, {0b0001101, 7}, {0b00010001, 8}, {0b00001000, 8}, {0b000000100, 9}},
	{{0b0001100, 7}, {0b0001011, 7}, {0b00010010, 8}, {0b000001111, 9}, {0b000001011, 9}, {0b000000010, 9}},
	{{0b0000111, 7}, {0b0000110, 7}, {0b00001001, 8}, {0b000001110, 9}, {0b000000011, 9}, {0b0000000001, 10}},
	{{0b00000110, 8}, {0b00000100, 8}, {0b000000101, 9}, {0b0000000011, 10}, {0b0000000010, 10}, {0b0000000000, 10}},
	// 8
	{{0b11, 2}, {0b100, 3}, {0b000110, 6}, {0b00010010, 8}, {0b00001100, 8}, {0b000000101, 9}},
	{{0b101, 3}, {0b01, 2}, {0b0010, 4}, {0b00010000, 8}, {0b00001001, 8}, {0b00000011, 8}},
	{{0b000111, 6}, {0b0011, 4}, {0b000101, 6}, {0b00001110, 8}, {0b00000111, 8}, {0b000000011, 9}},
	{{0b00010011, 8}, {0b00010001, 8}, {0b00001111, 8}, {0b000001101, 9}, {0b000001010, 9}, {0b0000000100, 10}},
	{{0b00001101, 8}, {0b0000101, 7}, {0b00001000, 8}, {0b000001011, 9}, {0b0000000101, 10}, {0b0000000001, 10}},
	{{0b000001100, 9}, {0b00000100, 8}, {0b000000100, 9}, {0b000000001, 9}, {0b00000000001, 11}, {0b00000000000, 11}},
	// 9
	{{0b111, 3}, {0b101, 3}, {0b01001, 5}, {0b001110, 6}, {0b00001111, 8}, {0b000000111, 9}},
	{{0b110, 3}, {0b100, 3}, {0b0101, 4}, {0b00101, 5}, {0b000110, 6}, {0b00000111, 8}},
	{{0b0111, 4}, {0b0110, 4}, {0b01000, 5}, {0b001000, 6}, {0b0001000, 7}, {0b00000101, 8}},
	{{0b001111, 6}, {0b00110, 5}, {0b001001, 6}, {0b0001010, 7}, {0b0000101, 7}, {0b00000001, 8}},
	{{0b0001011, 7}, {0b000111, 6}, {0b0001001, 7}, {0b0000110, 7}, {0b00000100, 8}, {0b000000001, 9}},
	{{0b00001110, 8}, {0b0000100, 7}, {0b00000110, 8}, {0b00000010, 8}, {0b000000110, 9}, {0b000000000, 9}},
	// 10
	{{0b1, 1}, {0b010, 3}, {0b001010, 6}, {0b00010111, 8}, {0b000100011, 9}, {0b000011110, 9}, {0b000001100, 9}, {0b0000010001, 10}},
	{{0b011, 3}, {0b0011, 4}, {0b001000, 6}, {0b0001100, 7}, {0b00010010, 8}, {0b000010101, 9}, {0b00001100, 8}, {0b00000111, 8}},
	{{0b001011, 6}, {0b001001, 6}, {0b0001111, 7}, {0b00010101, 8}, {0b000100000, 9}, {0b0000101000, 10}, {0b000010011, 9}, {0b000000110, 9}},
	{{0b0001110, 7}, {0b0001101, 7}, {0b00010110, 8}, {0b000100010, 9}, {0b0000101110, 10}, {0b0000010111, 10}, {0b000010010, 9}, {0b0000000111, 10}},
	{{0b00010100, 8}, {0b00010011, 8}, {0b000100001, 9}, {0b0000101111, 10}, {0b0000011011, 10}, {0b0000010110, 10}, {0b0000001001, 10}, {0b0000000011, 10}},
	{{0b000011111, 9}, {0b000010110, 9}, {0b0000101001, 10}, {0b0000011010, 10}, {0b00000010101, 11}, {0b00000010100, 11}, {0b0000000101, 10}, {0b00000000011, 11}},
	{{0b00001110, 8}, {0b00001101, 8}, {0b000001010, 9}, {0b0000001011, 10}, {0b0000010000, 10}, {0b0000000110, 10}, {0b00000000101, 11}, {0b00000000001, 11}},
	{{0b000001001, 9}, {0b00001000, 8}, {0b000000111, 9}, {0b0000001000, 10}, {0b0000000100, 10}, {0b00000000100, 11}, {0b00000000010, 11}, {0b00000000000, 11}},
	// 11
	{{0b11, 2}, {0b01010, 3}, {0b01010, 5}, {0b0011000, 7}, {0b00100010, 8}, {0b000100001, 9}, {0b00010101, 8}, {0b000001111, 9}},
	{{0b101, 3}, {0b011, 3}, {0b0100, 4}, {0b001010, 6}, {0b00100000, 8}, {0b00010001, 8}, {0b0001011, 7}, {0b00001010, 8}},
	{{0b01011, 5}, {0b00111, 5}, {0b001101, 6}, {0b0010010, 7}, {0b00011110, 8}, {0b000011111, 9}, {0b00010100, 8}, {0b00000101, 8}},
	{{0b0011001, 7}, {0b001011, 6}, {0b0010011, 7}, {0b000111011, 9}, {0b00011011, 8}, {0b0000010010, 10}, {0b00001100, 8}, {0b000000101, 9}},
	{{0b00100011, 8}, {0b00100001, 8}, {0b00011111, 8}, {0b000111010, 9}, {0b000011110, 9}, {0b0000010000, 10}, {0b000000111, 9}, {0b0000000101, 10}},
	{{0b00011100, 8}, {0b00011010, 8}, {0b000100000, 9}, {0b0000010011, 10}, {0b0000010001, 10}, {0b00000001111, 11}, {0b0000001000, 10}, {0b00000001110, 11}},
	{{0b00001110, 8}, {0b0001100, 7}, {0b0001001, 7}, {0b00001101, 8}, {0b000001110, 9}, {0b0000001001, 10}, {0b0000000100, 10}, {0b0000000001, 10}},
	{{0b00001011, 8}, {0b0000100, 7}, {0b00000110, 8}, {0b000000110, 9}, {0b0000000110, 10}, {0b0000000011, 10}, {0b0000000010, 10}, {0b0000000000, 10}},
	// 12
	{{0b1001, 4}, {0b110, 3}, {0b10000, 5}, {0b0100001, 7}, {0b00101001, 8}, {0b000100111, 9}, {0b000100110, 9}, {0b000011010, 9}},
	{{0b111, 3}, {0b101, 3}, {0b0110, 4}, {0b01001, 5}, {0b0010111, 7}, {0b0010000, 7}, {0b00011010, 8}, {0b00001011, 8}},
	{{0b10001, 5}, {0b0111, 4}, {0b01011, 5}, {0b001110, 6}, {0b0010101, 7}, {0b00011110, 8}, {0b0001010, 7}, {0b00000111, 8}},
	{{0b010001, 6}, {0b01010, 5}, {0b001111, 6}, {0b001100, 6}, {0b0010010, 7}, {0b00011100, 8}, {0b00001110, 8}, {0b00000101, 8}},
	{{0b0100000, 7}, {0b001101, 6}, {0b0010110, 7}, {0b0010011, 7}, {0b00010010, 8}, {0b00010000, 8}, {0b00001001, 8}, {0b000000101, 9}},
	{{0b00101000, 8}, {0b0010001, 7}, {0b00011111, 8}, {0b00011101, 8}, {0b00010001, 8}, {0b000001101, 8}, {0b00000100, 8}, {0b000000010, 9}},
	{{0b00011011, 8}, {0b0001100, 7}, {0b0001011, 7}, {0b00001111, 8}, {0b00001010, 8}, {0b000000111, 9}, {0b000000100, 9}, {0b0000000001, 10}},
	{{0b000011011, 9}, {0b00001100, 8}, {0b00001000, 8}, {0b000001100, 9}, {0b000000110, 9}, {0b000000011, 9}, {0b000000001, 9}, {0b0000000000, 10}},
	// 13
	{{0b1, 1}, {0b0101, 4}, {0b001110, 6}, {0b0010101, 7}, {0b00100010, 8}, {0b000110011, 9}, {0b000101110, 9}, {0b0001000111, 10},
		{0b000101010, 9}, {0b0000110100, 10}, {0b00001000100, 11}, {0b00000110100, 11}, {0b000001000011, 12}, {0b000000101100, 12}, {0b0000000101011, 13}, {0b0000000010011, 13}},
	{{0b011, 3}, {0b0100, 4}, {0b001100, 6}, {0b0010011, 7}, {0b00011111, 8}, {0b00011010, 8}, {0b000101100, 9}, {0b000100001, 9},
		{0b000011111, 9}, {0b000011000, 9}, {0b0000100000, 10}, {0b0000011000, 10}, {0b00000011111, 11}, {0b000000100011, 12}, {0b000000010110, 12}, {0b000000001110, 12}},
	{{0b001111, 6}, {0b001101, 6}, {0b0010111, 7}, {0b00100100, 8}, {0b000111011, 9}, {0b000110001, 9}, {0b0001001101, 10}, {0b0001000001, 10},
		{0b000011101, 9}, {0b0000101000, 10}, {0b0000011110, 11}, {0b00000101000, 11}, {0b00000011011, 11}, {0b000000100001, 12}, {0b0000000101010, 13}, {0b0000000010000, 13}},
	{{0b0010110, 7}, {0b0010100, 7}, {0b00100101, 8}, {0b000111101, 9}, {0b000111000, 9}, {0b0001001111, 10}, {0b0001001001, 10}, {0b0001000000, 10},
		{0b0000101011, 10}, {0b00001001100, 11}, {0b00000111000, 11}, {0b00000100101, 11}, {0b00000011010, 11}, {0b000000011111, 12}, {0b0000000011001, 13}, {0b0000000001110, 13}},
	{{0b00100011, 8}, {0b0010000, 7}, {0b000111100, 9}, {0b000111001, 9}, {0b0001100001, 10}, {0b0001100001, 10}, {0b00001110010, 11}, {0b00001011011, 11},
		{0b0000110110, 10}, {0b00001001001, 11}, {0b00000110111, 11}, {0b000000101001, 12}, {0b000000110000, 12}, {0b0000000110101, 13}, {0b0000000010111, 13}, {0b00000000011000, 14}},
	{{0b000111010, 9}, {0b00011011, 8}, {0b000110010, 9}, {0b0001100000, 10}, {0b0001001100, 10}, {0b0001000110, 10}, {0b00001011101, 11}, {0b00001010100, 11},
		{0b00001001101, 11}, {0b00000111010, 11}, {0b000001001111, 12}, {0b00000011101, 11}, {0b0000001001010, 13}, {0b0000000110001, 13}, {0b00000000101001, 14}, {0b00000000010001, 14}},
	{{0b000101111, 9}, {0b000101101, 9}, {0b0001001110, 10}, {0b0001001010, 10}, {0b00001110011, 11}, {0b00001011110, 11}, {0b00001011010, 11}, {0b00001001111, 11},
		{0b00001000101, 11}, {0b000001010011, 12}, {0b000001000111, 12}, {0b000000110010, 12}, {0b0000000111011, 13}, {0b0000000100110, 13}, {0b00000000100100, 14}, {0b00000000001111, 14}},
	{{0b0001001000, 10}, {0b000100010, 9}, {0b0000111000, 10}, {0b00001011111, 11}, {0b00001011100, 11}, {0b00001010101, 11}, {0b000001011011, 12}, {0b000001011010, 12},
		{0b000001010110, 12}, {0b000001001001, 12}, {0b0000001001101, 13}, {0b0000001000001, 13}, {0b0000000110011, 13}, {0b00000000101100, 14}, {0b0000000000101011, 16}, {0b0000000000101010, 16}},
	{{0b000101011, 9}, {0b00010100, 8}, {0b000011110, 9}, {0b0000101100, 10}, {0b0000110111, 10}, {0b00001001110, 11}, {0b00001001000, 11}, {0b000001010111, 12},
		{0b000001001110, 12}, {0b000000111101, 12}, {0b000000101110, 12}, {0b0000000110110, 13}, {0b0000000100101, 13}, {0b00000000011110, 14}, {0b000000000010100, 15}, {0b000000000010000, 15}},
	{{0b0000110101, 10}, {0b000011001, 9}, {0b0000101001, 10}, {0b0000100101, 10}, {0b00000101100, 11}, {0b00000111011, 11}, {0b00000110110, 11}, {0b0000001010001, 13},
		{0b000001000010, 12}, {0b0000001001100, 13}, {0b0000000111001, 13}, {0b00000000110110, 14}, {0b00000000100101, 14}, {0b00000000010010, 14}, {0b0000000000100111, 16}, {0b000000000001011, 15}},
	{{0b0000100011, 10}, {0b0000100001, 10}, {0b0000011111, 10}, {0b00000111001, 11}, {0b00000101010, 11}, {0b000001010010, 12}, {0b000001001000, 12}, {0b0000001010000, 13},
		{0b000000101111, 12}, {0b0000000111010, 13}, {0b00000000110111, 14}, {0b0000000010101, 13}, {0b00000000010110, 14}, {0b000000000011010, 15}, {0b0000000000100110, 16}, {0b00000000000010110, 17}},
	{{0b00000110101, 11}, {0b0000011001, 10}, {0b0000010111, 10}, {0b00000100110, 11}, {0b000001000110, 12}, {0b000000111100, 12}, {0b000000110011, 12}, {0b000000100100, 12},
		{0b0000000110111, 13}, {0b0000000011010, 13}, {0b0000000100010, 13}, {0b00000000010111, 14}, {0b000000000011011, 15}, {0b000000000001110, 15}, {0b000000000001001, 15}, {0b0000000000000111, 16}},
	{{0b00000100010, 11}, {0b00000100000, 11}, {0b00000011100, 11}, {0b000000100111, 12}, {0b000000110001, 12}, {0b0000001001011, 13}, {0b000000011110, 12}, {0b0000000110100, 13},
		{0b00000000110000, 14}, {0b00000000101000, 14}, {0b000000000110100, 15}, {0b000000000011100, 15}, {0b000000000010010, 15}, {0b0000000000010001, 16}, {0b0000000000001001, 16}, {0b0000000000000101, 16}},
	{{0b000000101101, 12}, {0b00000010101, 11}, {0b000000100010, 12}, {0b0000001000000, 13}, {0b0000000111000, 13}, {0b0000000110010, 13}, {0b00000000110001, 14}, {0b00000000101101, 14},
		{0b00000000011111, 14}, {0b00000000010011, 14}, {0b00000000001100, 14}, {0b000000000001111, 15}, {0b0000000000001010, 16}, {0b000000000000111, 15}, {0b0000000000000110, 16}, {0b0000000000000011, 16}},
	{{0b000000010000, 12}, {0b000000010111, 12}, {0b000000010100, 12}, {0b0000000100111, 13}, {0b0000000100100, 13}, {0b0000000100011, 13}, {0b000000000110101, 15}, {0b00000000010101, 14},
		{0b00000000010000, 14}, {0b00000000000010111, 17}, {0b000000000001101, 15}, {0b000000000001010, 15}, {0b000000000000110, 15}, {0b00000000000000001, 17}, {0b0000000000000100, 16}, {0b0000000000000010, 16}},
	{{0b000000010000, 12}, {0b000000001111, 12}, {0b0000000010001, 13}, {0b00000000011011, 14}, {0b00000000011001, 14}, {0b00000000010100, 14}, {0b000000000011101, 15}, {0b00000000001011, 14},
		{0b000000000010001, 15}, {0b000000000001100, 15}, {0b0000000000010000, 16}, {0b0000000000001000, 16}, {0b0000000000000000001, 19}, {0b000000000000000001, 18}, {0b0000000000000000000, 19}, {0b0000000000000001, 16}},
	// 14 Not used
	// 15
	{{0b111, 3}, {0b1100, 4}, {0b10010, 5}, {0b0110101, 7}, {0b0101111, 7}, {0b01001100, 8}, {0b001111100, 9}, {0b001101100, 9},
		{0b001011001, 9}, {0b0001111011, 10}, {0b0001101100, 10}, {0b00001110111, 11}, {0b00001101011, 11}, {0b00001010001, 11}, {0b000001111010, 12}, {0b0000000111111, 13}},
	{{0b1101, 4}, {0b101, 3}, {0b10000, 5}, {0b011011, 6}, {0b0101110, 7}, {0b0100100, 7}, {0b00111101, 8}, {0b00110011, 8},
		{0b00101010, 8}, {0b001000110, 9}, {0b000110100, 9}, {0b0001010011, 10}, {0b0001000001, 10}, {0b0000101001, 10}, {0b00000111011, 11}, {0b00000100100, 11}},
	{{0b10011, 5}, {0b10001, 5}, {0b01111, 5}, {0b0101001, 6}, {0b0101001, 7}, {0b0100010, 7}, {0b00111011, 8}, {0b00110000, 9},
		{0b00101000, 8}, {0b001000000, 9}, {0b000110010, 9}, {0b0001001110, 10}, {0b0000111110, 10}, {0b00001010000, 11}, {0b00000111000, 11}, {0b00000100001, 11}},
	{{0b011101, 6}, {0b011100, 6}, {0b011001, 6}, {0b0101011, 7}, {0b0100111, 7}, {0b00111111, 8}, {0b00110111, 8}, {0b001011101, 9},
		{0b001001100, 9}, {0b000111011, 9}, {0b0001011101, 10}, {0b0001001000, 10}, {0b0000110110, 10}, {0b00001001011, 11}, {0b00000110010, 11}, {0b00000011101, 11}},
	{{0b0110100, 7}, {0b010110, 6}, {0b0101010, 7}, {0b0101000, 7}, {0b01000011, 8}, {0b00111001, 8}, {0b001011111, 9}, {0b001001111, 9},
		{0b001001000, 9}, {0b000111001, 9}, {0b0001011001, 10}, {0b0001000101, 10}, {0b0000110001, 10}, {0b00001000010, 11}, {0b00000101110, 11}, {0b00000011011, 11}},
	{{0b01001101, 8}, {0b0100101, 7}, {0b0100011, 7}, {0b01000010, 8}, {0b00111010, 8}, {0b00110100, 8}, {0b001011011, 9}, {0b001001010, 9},
		{0b000111110, 9}, {0b000110000, 9}, {0b0001001111, 10}, {0b0000111111, 10}, {0b00001011010, 11}, {0b00000111110, 11}, {0b00000101000, 11}, {0b000000100110, 12}},
	{{0b001111101, 9}, {0b0100000, 7}, {0b00111100, 8}, {0b00111000, 8}, {0b00110010, 8}, {0b001011100, 9}, {0b001001110, 9}, {0b001000001, 9},
		{0b000110111, 9}, {0b0001010111, 10}, {0b0001000111, 10}, {0b0000110011, 10}, {0b00001001001, 11}, {0b00000110011, 11}, {0b000001000110, 12}, {0b000000011110, 12}},
	{{0b001101101, 9}, {0b00110101, 8}, {0b00110001, 8}, {0b001011110, 9}, {0b001011000, 9}, {0b001001011, 9}, {0b001000010, 9}, {0b0001111010, 10},
		{0b0001011011, 10}, {0b0001001001, 10}, {0b0000111000, 10}, {0b0000101010, 10}, {0b00001000000, 11}, {0b00000101100, 11}, {0b00000010101, 11}, {0b000000011001, 12}},
	{{0b001011010, 9}, {0b00101011, 8}, {0b00101001, 8}, {0b001001101, 9}, {0b001001001, 9}, {0b000111111, 9}, {0b000111000, 9}, {0b0001011100, 10},
		{0b0001001101, 10}, {0b0001000010, 10}, {0b0000101111, 10}, {0b00001000011, 11}, {0b00000110000, 11}, {0b000000110101, 12}, {0b000000100100, 12}, {0b000000010100, 12}},
	{{0b001000111, 9}, {0b00100010, 8}, {0b001000011, 9}, {0b000111100, 9}, {0b000111010, 9}, {0b000110001, 9}, {0b0001011000, 10}, {0b0001001100, 10},
		{0b0001000011, 10}, {0b00001101010, 11}, {0b00001000111, 11}, {0b00000110110, 11}, {0b00000100110, 11}, {0b000000100111, 12}, {0b000000010111, 12}, {0b000000001111, 12}},
	{{0b0001101101, 10}, {0b000110101, 9}, {0b000110011, 9}, {0b000101111, 9}, {0b0001011010, 10}, {0b0001010010, 10}, {0b0000111010, 10}, {0b0000111001, 10},
		{0b0000110000, 10}, {0b00001001000, 11}, {0b00000111001, 11}, {0b00000101001, 11}, {0b00000010111, 11}, {0b000000011011, 12}, {0b0000000111110, 13}, {0b000000001001, 12}},
	{{0b0001010110, 10}, {0b000101010, 9}, {0b000101000, 9}, {0b000100101, 9}, {0b0001000110, 10}, {0b0001000000, 10}, {0b0000110100, 10}, {0b0000101011, 10},
		{0b00001000110, 11}, {0b00000110111, 11}, {0b00000101010, 11}, {0b00000011001, 11}, {0b000000011101, 12}, {0b000000010010, 12}, {0b000000001011, 12}, {0b0000000001011, 13}},
	{{0b00001110110, 11}, {0b0001000100, 10}, {0b000011110, 9}, {0b0000110111, 10}, {0b0000110010, 10}, {0b0000101110, 10}, {0b00001001010, 11}, {0b00001000001, 11},
		{0b00000110001, 11}, {0b00000100111, 11}, {0b00000011000, 11}, {0b00000010000, 11}, {0b000000010110, 12}, {0b000000001101, 12}, {0b0000000001110, 13}, {0b0000000000111, 13}},
	{{0b00001011011, 11}, {0b0000101100, 10}, {0b0000100111, 10}, {0b0000100110, 10}, {0b0000100010, 10}, {0b00000111111, 11}, {0b00000110100, 11}, {0b00000101101, 11},
		{0b00000011111, 11}, {0b000000110100, 12}, {0b000000011100, 12}, {0b000000010011, 12}, {0b000000001110, 12}, {0b000000001000, 12}, {0b0000000001001, 13}, {0b0000000000011, 13}},
	{{0b000001111011, 12}, {0b00000111100, 11}, {0b00000111010, 11}, {0b00000110101, 11}, {0b00000101111, 11}, {0b00000101011, 11}, {0b00000100000, 11}, {0b00000010110, 11},
		{0b000000100101, 12}, {0b000000011000, 12}, {0b000000010001, 12}, {0b000000001100, 12}, {0b0000000001111, 13}, {0b0000000001010, 13}, {0b000000000010, 12}, {0b0000000000001, 13}},
	{{0b000001000111, 12}, {0b00000100101, 11}, {0b00000100010, 11}, {0b00000011110, 11}, {0b00000011100, 11}, {0b00000010100, 11}, {0b00000010001, 11}, {0b000000011010, 12},
		{0b000000010101, 12}, {0b000000010000, 12}, {0b000000001010, 12}, {0b000000000110, 12}, {0b0000000001000, 13}, {0b0000000000110, 13}, {0b0000000000010, 13}, {0b0000000000000, 13}},
	// 16 17 18 19 20 21 22 23
	{{0b1, 1}, {0b0101, 4}, {0b001110, 6}, {0b00101100, 8}, {0b001001010, 9}, {0b000111111, 9}, {0b0001101110, 10}, {0b0001011101, 10},
		{0b00010101100, 11}, {0b00010010101, 11}, {0b00010001010, 11}, {0b000011110010, 12}, {0b000011100001, 12}, {0b000011000011, 12}, {0b0000101111000, 13}, {0b000010001, 9}},
	{{0b011, 3}, {0b0100, 4}, {0b001100, 6}, {0b0010100, 7}, {0b00100011, 8}, {0b000111110, 9}, {0b000110101, 9}, {0b000101111, 9},
		{0b0001010011, 10}, {0b0001001011, 10}, {0b0001000100, 10}, {0b00001110111, 11}, {0b000011001001, 12}, {0b00001101011, 11}, {0b000011001111, 12}, {0b00001001, 8}},
	{{0b001111, 6}, {0b001101, 6}, {0b0010111, 7}, {0b00100110, 8}, {0b001000011, 9}, {0b000111010, 9}, {0b0001100111, 10}, {0b0001011010, 10},
		{0b00010100001, 11}, {0b0001001000, 10}, {0b00001111111, 11}, {0b00001110101, 11}, {0b00001101110, 11}, {0b000011010001, 12}, {0b000011001110, 12}, {0b000010000, 9}},
	{{0b00101101, 8}, {0b0010101, 7}, {0b00100111, 8}, {0b001000101, 9}, {0b001000000, 9}, {0b0001110010, 10}, {0b0001100011, 10}, {0b0001010111, 10},
		{0b00010011110, 11}, {0b00010001100, 11}, {0b000011111100, 12}, {0b000011010100, 12}, {0b000011000111, 12}, {0b0000110000011, 13}, {0b0000101101101, 13}, {0b0000011010, 10}},
	{{0b001001011, 9}, {0b00100100, 8}, {0b001000100, 9}, {0b001000001, 9}, {0b0001110011, 10}, {0b0001100101, 10}, {0b00010110011, 11}, {0b00010100100, 11},
		{0b00010011011, 11}, {0b000100001000, 12}, {0b000011110110, 12}, {0b000011100010, 12}, {0b0000110001011, 13}, {0b0000101111110, 13}, {0b0000101101010, 13}, {0b000001001, 9}},
	{{0b001000010, 9}, {0b00011110, 8}, {0b000111011, 9}, {0b000111000, 9}, {0b0001100110, 10}, {0b00010111001, 11}, {0b00010101101, 11}, {0b000100001001, 12},
		{0b00010001110, 11}, {0b000011111101, 12}, {0b000011101000, 12}, {0b0000110010000, 13}, {0b0000110000100, 13}, {0b0000101111010, 13}, {0b00000110111101, 14}, {0b0000010000, 10}},
	{{0b0001101111, 10}, {0b000110110, 9}, {0b000110100, 9}, {0b0001100100, 10}, {0b00010111000, 11}, {0b00010110010, 11}, {0b00010100000, 11}, {0b00010000101, 11},
		{0b000100000001, 12}, {0b000011110100, 12}, {0b000011100100, 12}, {0b000011011001, 12}, {0b0000110000001, 13}, {0b0000101101110, 13}, {0b00001011001011, 14}, {0b0000001010, 10}},
	{{0b0001100010, 10}, {0b000110000, 9}, {0b0001011011, 10}, {0b0001011000, 10}, {0b00010100101, 11}, {0b00010011101, 11}, {0b00010010100, 11}, {0b000100000101, 12},
		{0b000011111000, 12}, {0b0000110010111, 13}, {0b0000110001101, 13}, {0b0000101110100, 13}, {0b0000101111100, 13}, {0b000001101111001, 15}, {0b000001101110100, 15}, {0b0000001000, 10}},
	{{0b0001010101, 10}, {0b0001010100, 10}, {0b0001010001, 10}, {0b00010011111, 11}, {0b00010011100, 11}, {0b00010001111, 11}, {0b000100000100, 12}, {0b000011111001, 12},
		{0b0000110101011, 13}, {0b0000110010001, 13}, {0b0000110001000, 13}, {0b0000101111111, 13}, {0b00001011010111, 14}, {0b00001011001001, 14}, {0b00001011000100, 14}, {0b0000000111, 10}},
	{{0b00010011010, 11}, {0b0001001100, 10}, {0b0001001001, 10}, {0b00010001101, 11}, {0b00010000011, 11}, {0b000100000000, 12}, {0b000011110101, 12}, {0b0000110101010, 13},
		{0b0000110010110, 13}, {0b0000110001010, 13}, {0b0000110000000, 13}, {0b00001011011111, 14}, {0b0000101100111, 13}, {0b00001011000110, 14}, {0b0000101100000, 13}, {0b00000001011, 11}},
	{{0b00010001011, 11}, {0b00010000001, 11}, {0b0001000011, 10}, {0b00001111101, 11}, {0b000011110111, 12}, {0b000011101001, 12}, {0b000011100101, 12}, {0b000011011011, 12},
		{0b0000110001001, 13}, {0b00001011100111, 14}, {0b00001011100001, 14}, {0b00001011010000, 14}, {0b000001101110101, 15}, {0b000001101110010, 15}, {0b00000110110111, 14}, {0b0000000100, 10}},
	{{0b000011110011, 12}, {0b00001111000, 11}, {0b00001110110, 11}, {0b00001110011, 11}, {0b000011100011, 12}, {0b000011011111, 12}, {0b0000110001100, 13}, {0b00001011101010, 14},
		{0b00001011100110, 14}, {0b00001011100000, 14}, {0b00001011010001, 14}, {0b00001011001000, 14}, {0b00001011000010, 14}, {0b0000011011111, 13}, {0b00000110110100, 14}, {0b00000000110, 11}},
	{{0b000011001010, 12}, {0b000011100000, 12}, {0b000011011110, 12}, {0b000011011010, 12}, {0b000011011000, 12}, {0b0000110000101, 13}, {0b0000110000010, 13}, {0b0000101111101, 13},
		{0b0000101101100, 13}, {0b000001101111000, 15}, {0b00000110111011, 14}, {0b00001011000011, 14}, {0b00000110111000, 14}, {0b00000110110101, 14}, {0b0000011011000000, 16}, {0b00000000100, 11}},
	{{0b00001011101011, 14}, {0b000011010011, 12}, {0b000011010010, 12}, {0b000011010000, 12}, {0b0000101110010, 13}, {0b0000101111011, 13}, {0b00001011011110, 14}, {0b00001011010011, 14},
		{0b00001011001010, 14}, {0b0000011011000111, 16}, {0b000001101110011, 15}, {0b000001101101101, 15}, {0b000001101101100, 15}, {0b00000110110000011, 17}, {0b000001101100001, 15}, {0b00000000010, 11}},
	{{0b0000101111001, 13}, {0b0000101110001, 13}, {0b00001100110, 11}, {0b000010111011, 12}, {0b00001011010110, 14}, {0b00001011010010, 14}, {0b0000101100110, 13}, {0b00001011000111, 14},
		{0b00001011000101, 14}, {0b000001101100010, 15}, {0b0000011011000110, 16}, {0b000001101100111, 15}, {0b00000110110000010, 17}, {0b000001101100110, 15}, {0b00000110110010, 14}, {0b00000000000, 11}},
	{{0b000001100, 9}, {0b00001010, 8}, {0b00000111, 8}, {0b000001011, 9}, {0b000001010, 9}, {0b0000010001, 10}, {0b0000001011, 10}, {0b0000001001, 10},
		{0b00000001101, 11}, {0b00000001100, 11}, {0b00000001010, 11}, {0b00000000111, 11}, {0b00000000101, 11}, {0b00000000011, 11}, {0b00000000001, 11}, {0b00000011, 8}},
	// 24 25 26 27 28 29 30 31
	{{0b1111, 4}, {0b1101, 4}, {0b101110, 6}, {0b1010000, 7}, {0b10010010, 8}, {0b100000110, 9}, {0b011111000, 9}, {0b0110110010, 10},
		{0b0110101010, 10}, {0b01010011101, 11}, {0b01010001101, 11}, {0b01010001001, 11}, {0b01001101101, 11}, {0b01000000101, 11}, {0b010000001000, 12}, {0b001011000, 9}},
	{{0b1110, 4}, {0b1100, 4}, {0b10101, 5}, {0b100110, 6}, {0b1000111, 7}, {0b10000010, 8}, {0b01111010, 8}, {0b011011000, 9},
		{0b011010001, 9}, {0b011000110, 9}, {0b0101000111, 10}, {0b0101011001, 10}, {0b0100111111, 10}, {0b0100101001, 10}, {0b0100010111, 10}, {0b00101010, 8}},
	{{0b101111, 6}, {0b10110, 5}, {0b101001, 6}, {0b1001010, 7}, {0b1000100, 7}, {0b10000000, 8}, {0b01111000, 8}, {0b011011101, 9},
		{0b011001111, 9}, {0b011000010, 9}, {0b010110110, 9}, {0b0101010100, 10}, {0b0100111011, 10}, {0b0100100111, 11}, {0b01000011101, 11}, {0b0010010, 7}},
	{{0b1010001, 7}, {0b100111, 6}, {0b1001011, 7}, {0b1000110, 7}, {0b10000110, 8}, {0b01111101, 8}, {0b01110100, 8}, {0b011011100, 9},
		{0b011001100, 9}, {0b010111110, 9}, {0b010110010, 9}, {0b0101000101, 10}, {0b0100110111, 10}, {0b0100100101, 10}, {0b0100001111, 10}, {0b0010000, 7}},
	{{0b10010011, 8}, {0b1001000, 7}, {0b1000101, 7}, {0b10000111, 8}, {0b01111111, 8}, {0b01110110, 8}, {0b01110000, 8}, {0b011010010, 9},
		{0b011001000, 9}, {0b010111100, 9}, {0b0101100000, 10}, {0b0101000011, 10}, {0b0100110010, 10}, {0b0100011101, 10}, {0b01000011100, 11}, {0b0001110, 7}},
	{{0b100000111, 9}, {0b1000010, 7}, {0b10000001, 8}, {0b01111110, 8}, {0b01110111, 8}, {0b01110010, 8}, {0b011010110, 9}, {0b011001010, 9},
		{0b011000000, 9}, {0b010110100, 9}, {0b0101010101, 10}, {0b0100111101, 10}, {0b0100101101, 10}, {0b0100011001, 10}, {0b0100000110, 10}, {0b0001100, 7}},
	{{0b011111001, 9}, {0b01111011, 8}, {0b01111001, 8}, {0b01110101, 8}, {0b01110001, 8}, {0b011010111, 9}, {0b011001110, 9}, {0b011000011, 9},
		{0b010111001, 9}, {0b0101011011, 10}, {0b0101001010, 10}, {0b0100110100, 10}, {0b0100100011, 10}, {0b0100010000, 10}, {0b01000001000, 11}, {0b0001010, 7}},
	{{0b0110110011, 10}, {0b01110011, 8}, {0b01101111, 8}, {0b01101101, 8}, {0b011010011, 9}, {0b011001011, 9}, {0b011000100, 9}, {0b010111011, 9},
		{0b0101100001, 10}, {0b0101001100, 10}, {0b0100111001, 10}, {0b0100101010, 10}, {0b0100011011, 10}, {0b01000010011, 11}, {0b00101111101, 11}, {0b00010001, 8}},
	{{0b0110101011, 10}, {0b011010100, 9}, {0b011010000, 9}, {0b011001101, 9}, {0b011001001, 9}, {0b011000001, 9}, {0b010111010, 9}, {0b010110001, 9},
		{0b010101001, 9}, {0b0101000000, 10}, {0b0100101111, 10}, {0b0100011110, 10}, {0b0100001100, 10}, {0b01000000010, 11}, {0b00101111001, 11}, {0b00010000, 8}},
	{{0b0101001111, 10}, {0b011000111, 9}, {0b011000101, 9}, {0b010111111, 9}, {0b010111101, 9}, {0b010110101, 9}, {0b010101110, 9}, {0b0101001101, 10},
		{0b0101000001, 10}, {0b0100110001, 10}, {0b0100100001, 10}, {0b0100010011, 10}, {0b01000001001, 11}, {0b00101111011, 11}, {0b00101110011, 11}, {0b00001011, 8}},
	{{0b01010011100, 11}, {0b010111000, 9}, {0b010110111, 9}, {0b010110011, 9}, {0b010101111, 9}, {0b0101011000, 10}, {0b0101001011, 10}, {0b0100111010, 10},
		{0b0100110000, 10}, {0b0100100010, 10}, {0b0100010101, 10}, {0b01000010010, 11}, {0b00101111111, 11}, {0b00101110101, 11}, {0b00101101110, 11}, {0b00001010, 8}},
	{{0b01010001100, 11}, {0b0101011010, 10}, {0b010101011, 9}, {0b010101000, 9}, {0b010100100, 9}, {0b0100111110, 10}, {0b0100110101, 10}, {0b0100101011, 10},
		{0b0100011111, 10}, {0b0100010100, 10}, {0b0100000111, 10}, {0b01000000001, 11}, {0b00101110111, 11}, {0b00101110000, 11}, {0b00101101010, 11}, {0b00000110, 8}},
	{{0b01010001000, 11}, {0b0101000010, 10}, {0b0100111100, 10}, {0b0100111000, 10}, {0b0100110011, 10}, {0b0100101110, 10}, {0b0100100100, 10}, {0b0100011100, 10},
		{0b0100001101, 10}, {0b0100000101, 10}, {0b01000000000, 11}, {0b00101111000, 11}, {0b00101110010, 11}, {0b00101101100, 11}, {0b00101100111, 11}, {0b00000100, 8}},
	{{0b01001101100, 11}, {0b0100101100, 10}, {0b0100101000, 10}, {0b0100100110, 10}, {0b0100100000, 10}, {0b0100011010, 10}, {0b0100010001, 10}, {0b0100001010, 10},
		{0b01000000011, 11}, {0b00101111100, 11}, {0b00101110110, 11}, {0b00101110001, 11}, {0b00101101101, 11}, {0b00101101001, 11}, {0b00101100101, 11}, {0b00000010, 8}},
	{{0b010000001001, 12}, {0b0100011000, 10}, {0b0100010110, 10}, {0b0100010010, 10}, {0b0100001011, 10}, {0b0100001000, 10}, {0b0100000011, 10}, {0b00101111110, 11},
		{0b00101111010, 11}, {0b00101110100, 11}, {0b00101101111, 11}, {0b00101101011, 11}, {0b00101101000, 11}, {0b00101100110, 11}, {0b00101100100, 8}, {0b00000000, 8}},
	{{0b00101011, 8}, {0b0010100, 7}, {0b0010011, 7}, {0b0010001, 7}, {0b0001111, 7}, {0b0001101, 7}, {0b0001011, 7}, {0b0001001, 7},
		{0b0000111, 7}, {0b0000110, 7}, {0b0000100, 7}, {0b00000111, 8}, {0b00000101, 8}, {0b00000011, 8}, {0b00000001, 8}, {0b0011, 4}},
}

var huffmanTableA = [16][3]int{
	{0b1, 1, 0000}, {0b0101, 4, 0001}, {0b0100, 4, 0010}, {0b00101, 5, 0011},
	{0b0110, 4, 0100}, {0b000101, 6, 0101}, {0b00100, 5, 0110}, {0b000100, 6, 0111},
	{0b011, 4, 1000}, {0b00011, 5, 1001}, {0b00110, 5, 1010}, {0b000000, 6, 1011},
	{0b00111, 5, 1100}, {0b000010, 6, 1101}, {0b000011, 6, 1110}, {0b000001, 6, 1111},
}

var huffmanTableB = [16][3]int{
	{0b1111, 4, 0000}, {0b1110, 4, 0001}, {0b1101, 4, 0010}, {0b1100, 4, 0011},
	{0b1011, 4, 0100}, {0b1010, 4, 0101}, {0b1001, 4, 0110}, {0b1000, 4, 0111},
	{0b0111, 4, 1000}, {0b0110, 4, 1001}, {0b0101, 4, 1010}, {0b0100, 4, 1011},
	{0b0011, 4, 1100}, {0b0010, 4, 1101}, {0b0001, 4, 1110}, {0b0000, 4, 1111},
}

type huffmanTable struct {
	Table   [][][2]int
	Linbits int
}

var huffmanTables = [...]huffmanTable{
	{nil, 0},                    // Table 0 Not used
	{huffmanCodes[0:2], 0},      // Table 1
	{huffmanCodes[2:5], 0},      // Table 2
	{huffmanCodes[5:8], 0},      // Table 3
	{nil, 0},                    // Table 4 Not used
	{huffmanCodes[8:12], 0},     // Table 5
	{huffmanCodes[12:16], 0},    // Table 6
	{huffmanCodes[16:22], 0},    // Table 7
	{huffmanCodes[22:28], 0},    // Table 8
	{huffmanCodes[28:34], 0},    // Table 9
	{huffmanCodes[34:42], 0},    // Table 10
	{huffmanCodes[42:50], 0},    // Table 11
	{huffmanCodes[50:58], 0},    // Table 12
	{huffmanCodes[58:74], 0},    // Table 13
	{nil, 0},                    // Table 14 Not used
	{huffmanCodes[74:90], 0},    // Table 15
	{huffmanCodes[90:106], 1},   // Table 16
	{huffmanCodes[90:106], 2},   // Table 17
	{huffmanCodes[90:106], 3},   // Table 18
	{huffmanCodes[90:106], 4},   // Table 19
	{huffmanCodes[90:106], 6},   // Table 20
	{huffmanCodes[90:106], 8},   // Table 21
	{huffmanCodes[90:106], 10},  // Table 22
	{huffmanCodes[90:106], 13},  // Table 23
	{huffmanCodes[106:122], 4},  // Table 24
	{huffmanCodes[106:122], 5},  // Table 25
	{huffmanCodes[106:122], 6},  // Table 26
	{huffmanCodes[106:122], 7},  // Table 27
	{huffmanCodes[106:122], 8},  // Table 28
	{huffmanCodes[106:122], 9},  // Table 29
	{huffmanCodes[106:122], 11}, // Table 30
	{huffmanCodes[106:122], 13}, // Table 31
}

func decodeHuffman(r *BitReader, tableNumber int) (x, y, v, w int) {
	table := huffmanTables[tableNumber]

	bitSample := r.ReadBits(24)
	for x, v := range table.Table {
		for y, k := range v {
			hcod := k[0]
			hlen := k[1]

			if hcod == bitSample>>(24-hlen) {
				r.Seek(-(24 - hlen))

				if x == maxTableEntry && table.Linbits > 0 {
					x += r.ReadBits(table.Linbits)
				}
				if x != 0 && r.ReadBits(1) == 1 {
					x = -x
				}

				if y == maxTableEntry && table.Linbits > 0 {
					y += r.ReadBits(table.Linbits) //
				}
				if y != 0 && r.ReadBits(1) == 1 {
					y = -y
				}

				return x, y, 0, 0
			}
		}
	}
	r.Seek(-24)
	return x, y, v, w
}

func decodeHuffmanB(r *BitReader) (v, w, x, y int) {
	bitSample := r.ReadBits(24)
	for _, k := range huffmanTableA {
		hcod := k[0]
		hlen := k[1]

		if hcod == bitSample>>(24-hlen) {
			r.Seek(-(24 - hlen))

			v = k[2] & 0x8
			w = k[2] & 0x4
			x = k[2] & 0x2
			y = k[2] & 0x1

			return v, w, x, y
		}
	}
	r.Seek(-24)
	return v, w, x, y
}