package tools

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"rasp-cloud/message"

	"rasp-cloud/config"

	"google.golang.org/protobuf/proto"
)

var normalMap [36]byte = [36]byte{'p', 'e', '5', 'b', 'm', 'o', '0', 'd', '3', 'l', '1', 'r', 'k', 'c', 'q', 'x', 'y', 'g', 'a', 'z', 'u', 'n', '7', 'j', '8', 's', '9', '4', 'v', 'i', 'h', '2', 'f', 'w', 't', '6'}
var idMapper [36][36]byte = [36][36]byte{{'p', 'e', '5', 'b', 'm', 'o', '0', 'd', '3', 'l', '1', 'r', 'k', 'c', 'q', 'x', 'y', 'g', 'a', 'z', 'u', 'n', '7', 'j', '8', 's', '9', '4', 'v', 'i', 'h', '2', 'f', 'w', 't', '6'}, {'i', '1', '9', 'x', '7', 'm', 'u', 'j', 'c', 'e', 'd', '4', 'n', 'z', '3', 'a', 't', 'w', '0', 'b', 'r', '8', '5', 'g', 'q', 's', 'y', 'l', 'p', 'v', 'o', '6', '2', 'k', 'f', 'h'}, {'3', 'r', '6', 'k', 'z', 'c', '2', 'v', 't', '9', 'n', 'h', '5', 'y', 'w', '8', 'o', 'x', 'b', 'l', '7', '1', 'i', 'p', 'g', 'f', 's', 'q', 'm', '0', 'd', 'u', '4', 'a', 'j', 'e'}, {'t', 'r', 's', 'e', '9', 'b', '0', '5', '1', '8', 'g', 'j', 'l', 'a', '6', 'w', '3', 'p', 'k', '4', 'z', 'u', 'o', 'c', 'v', '2', 'n', 'y', 'x', 'q', 'h', '7', 'f', 'd', 'm', 'i'}, {'k', '1', 'q', 'u', '3', 'y', '8', 'b', 'j', 'm', 'f', '7', '2', 't', 'a', 'p', 'r', 'd', 'w', 's', 'o', 'z', 'x', '6', 'h', 'g', 'v', '4', 'c', 'n', '9', 'e', 'l', '5', '0', 'i'}, {'m', 'u', 'b', '5', 's', '9', 'x', 'v', 'i', 'f', 'r', '6', 'k', '2', '4', '3', 'q', '0', 'g', 'h', 'c', 't', 'd', '1', 'a', 'z', 'e', 'o', 'j', 'w', 'p', 'n', 'l', '7', 'y', '8'}, {'v', 'r', 'k', 'x', 'i', 'u', 'y', '9', 'm', 'c', '5', 'j', '7', 'w', '6', 's', 'z', 'h', '0', 'e', '1', '2', 'g', '3', 'o', 'a', 'l', 'b', '8', '4', 'd', 'f', 'q', 'n', 'p', 't'}, {'z', 'v', 'b', 'p', 'm', 't', 's', 'n', 'h', 'o', 'w', '8', 'r', 'x', '0', 'j', 'l', 'q', 'k', '9', 'e', 'd', 'a', 'g', '5', '2', '7', 'i', '1', '4', '6', 'u', 'f', '3', 'c', 'y'}, {'l', 'x', 'g', 'q', 'v', 'b', '9', 'p', 'y', '8', 's', 'r', 'j', 'd', 't', '2', '1', 'w', 'a', 'f', 'n', 'o', '7', '3', 'u', '6', 'e', '5', 'i', '0', 'z', 'c', '4', 'm', 'k', 'h'}, {'8', 'r', 'f', 'w', 'a', 'n', 'h', 'g', '2', '4', 's', '6', 'v', '1', 'e', 'c', '5', '7', 'u', 'x', '9', 'y', 'k', '0', 'l', 'o', 'q', 'i', 'z', 'd', '3', 'j', 't', 'b', 'p', 'm'}, {'8', 'u', 'v', '3', 'e', 'z', 'q', 'c', '1', 'f', '0', 'x', 'g', '2', 'i', '9', 'y', 'p', 'k', 'n', 'o', '7', 'h', 'd', 'w', '6', '4', 'b', 'j', 'r', 's', 't', 'm', 'l', '5', 'a'}, {'u', 'd', 'x', 'y', '9', 'b', 'o', '2', '6', '3', 'j', 'l', 'n', 'i', 'z', 't', '8', '0', 's', '4', 'm', '7', 'f', 'a', '1', 'k', '5', 'q', 'p', 'v', 'r', 'e', 'h', 'w', 'g', 'c'}, {'p', 'x', 'j', 'y', '6', '4', 'a', 'q', 'r', '8', '5', '2', 'b', 'z', 's', '9', 'n', '1', '7', 't', '0', 'h', 'f', 'u', 'd', 'w', 'e', 'o', '3', 'i', 'c', 'g', 'k', 'v', 'l', 'm'}, {'x', 'v', 'p', 'a', 'e', 'd', 's', 'o', 'c', 'r', '6', 'i', 't', 'y', 'm', 'n', '3', '1', 'q', 'j', '5', '8', '9', 'u', 'k', 'w', 'g', 'z', 'f', '2', '7', '0', 'h', '4', 'l', 'b'}, {'x', '8', 'f', '7', 'p', 'g', 'm', '0', 'u', 'h', 'd', '2', 'e', 'o', 's', 't', 'r', 'v', 'q', '9', 'i', 'a', 'z', 'k', 'y', '4', '5', 'c', 'b', 'n', '6', 'j', 'l', '1', '3', 'w'}, {'e', 'x', 'f', 'k', 'z', 'o', 'q', 'l', '5', 'n', 'j', '7', '1', 'm', 'i', '8', 'u', 'y', 'c', '6', 'r', '3', 's', '2', '9', 'v', 'p', 'b', '4', 't', 'g', 'h', 'a', '0', 'd', 'w'}, {'8', 'h', 'u', 't', 'q', 'f', 'l', '0', '4', '9', '7', 'o', 'j', 'z', 'i', 'b', 'r', 'w', '5', 'c', 's', 'p', '6', 'n', 'e', '1', 'v', 'x', 'y', 'd', 'k', '3', 'm', '2', 'a', 'g'}, {'8', '7', 'b', '9', 'h', 'e', 'u', 't', '1', 'k', '0', '6', 'r', 's', 'n', 'm', 'o', '4', 'q', 'l', 'g', '5', 'p', 'x', 'j', 'i', 'z', '3', '2', 'c', 'a', 'd', 'y', 'w', 'f', 'v'}, {'m', 'b', 'n', 'z', 't', 'x', '0', '9', 'j', '5', 's', 'h', 'r', 'u', 'l', 'p', 'q', 'a', '1', 'i', 'o', '2', 'c', '6', '7', '8', 'e', 'd', '4', 'w', 'f', 'v', 'y', 'g', '3', 'k'}, {'e', 't', 'p', 'k', 'h', 'd', 'w', 'f', '0', 'x', '9', 'j', 'z', 'm', '2', 's', '7', 'u', 'v', 'l', 'n', '3', 'y', '1', '4', '8', 'b', 'c', '6', 'q', 'i', 'a', 'g', '5', 'r', 'o'}, {'0', 'n', 'd', 's', '9', 'k', '2', 'w', 'f', '5', 'x', 'z', 'y', 'h', 'b', '4', '7', 'o', 'a', 'u', 'r', 't', '1', 'm', 'j', 'l', 'i', '3', '6', 'g', 'e', 'c', '8', 'v', 'q', 'p'}, {'p', 'y', 'q', '7', 'w', 'l', '0', 'x', '4', 'r', 'k', 'd', '3', 'i', 'h', 'j', 'v', 'a', 't', 'z', 'c', 'b', 's', '1', 'u', '2', 'm', '5', 'e', '9', 'o', '8', 'f', 'g', 'n', '6'}, {'c', 'j', 'l', 'z', 'n', 'b', 'g', 'd', 'e', '5', 'q', '0', 'y', 'o', 'k', '9', 'u', 'a', 'r', 'v', 'w', 'i', '8', 'f', '2', 'm', 'h', 't', '7', '1', '3', '6', '4', 'p', 's', 'x'}, {'6', 'n', 'r', 'e', 'o', 'z', 'a', 'x', 'v', 's', '8', 'd', 'w', '9', '7', '2', '3', '1', 'b', '0', 'i', 'l', 'k', 'q', 'y', 'h', 'f', 'p', '4', 'c', 'u', 'j', 'g', 't', '5', 'm'}, {'c', 'k', 'd', 'l', 'g', 'i', 'u', '0', '5', 'h', '7', 'z', 'n', 't', 'w', '8', 's', '3', 'r', '1', 'p', 'f', 'j', 'a', 'q', '6', '2', '4', 'o', 'm', 'x', 'v', 'e', '9', 'y', 'b'}, {'p', 'f', 't', '6', 'l', 'b', 'd', '1', 'x', 'q', 'm', 'n', 'u', '7', 'a', '4', 'y', 'j', 's', '3', 'v', 'h', 'e', 'i', 'c', '5', '8', 'z', 'o', '2', '0', '9', 'r', 'k', 'w', 'g'}, {'b', 'd', 'p', 'o', '3', '4', 'w', 'l', 'u', 'c', 'v', '9', 'k', '0', '5', 'x', 'j', 's', '6', '7', 'y', 'n', '2', '8', 'e', 'f', 'a', 'g', 'r', '1', 'z', 'i', 't', 'm', 'q', 'h'}, {'n', '9', '5', 'u', 'g', 'z', '3', 'o', '7', 'k', 'f', 'v', '2', 't', 'p', 'q', '4', 'c', 'y', '6', 'w', 'h', 'r', 'x', 'j', 'b', '0', 'a', '8', 's', 'l', 'e', 'd', '1', 'i', 'm'}, {'p', 't', 's', 'u', 'g', 'w', 'z', '1', '3', '0', '2', 'y', '7', 'o', 'x', 'b', '8', '9', 'd', 'r', 'f', 'c', 'a', 'k', 'v', 'l', 'q', 'e', 'n', 'h', 'm', '4', 'j', '5', '6', 'i'}, {'y', '4', '8', 'k', 'x', 'e', 'a', 'b', 'p', 'z', 'i', '6', 'c', 'l', 's', 'v', '3', 'f', 'm', '0', 'n', '9', '2', 'g', '5', 'd', 'h', 'q', 'j', '1', 'r', 't', 'u', '7', 'o', 'w'}, {'s', 'c', '7', 'v', '2', 'n', 'z', 'b', 'p', 'h', 'r', 'u', '9', 'l', 't', 'k', 'w', '0', 'a', 'e', 'd', 'y', 'm', 'g', 'i', 'j', 'x', '3', '8', '4', '1', 'o', 'q', 'f', '6', '5'}, {'q', '3', 'j', 'y', 'u', '1', 'm', 'r', 'k', 'i', 'a', 'l', 'o', 'x', 'c', 'n', '2', 'g', 'b', 't', 'v', 'z', '4', '7', 'd', '8', '9', 'w', 'h', 's', '0', 'e', 'p', 'f', '6', '5'}, {'k', '7', 'o', 'v', 'b', 't', '0', '6', 'r', 'c', 'z', 'd', 'e', 'l', '9', 'u', 'q', 'f', '8', '5', '2', 'p', 'm', 'x', 's', 'a', 'n', '3', '1', 'i', 'h', 'y', '4', 'j', 'g', 'w'}, {'q', '4', 'f', 'c', '8', '7', 'j', 's', 'o', '2', 'v', 'l', 'g', 'm', 'n', 'b', '5', 'z', 'r', 't', '1', 'a', '9', 'k', '6', 'w', 'x', 'p', 'i', 'h', 'e', '3', '0', 'd', 'u', 'y'}, {'2', 'w', 'u', 'n', 'c', 'm', 'o', 'k', 'a', 'r', '6', 'x', 'b', 'y', 'p', 'z', 'i', 'l', '3', 't', 'f', '4', 'v', 'd', 'e', '1', '0', 'g', 's', '9', '8', '7', 'q', '5', 'j', 'h'}, {'4', '3', 'v', '9', 'k', 'n', 't', 'i', 'm', 'c', '0', 'f', '8', 'y', 's', 'q', 'e', 'w', 'r', '1', 'j', 'd', 'p', '6', 'h', '2', '5', 'g', 'o', 'a', 'b', '7', 'u', 'z', 'x', 'l'}}

// var normalMap [36]byte = [36]byte{'6', 'p', 'e', '5', 'b', 'm', 'o', '0', 'd', '3', 'l', '1', 'r', 'k', 'c', 'q', 'x', 'y', 'g', 'a', 'z', 'u', 'n', '7', 'j', '8', 's', '9', '4', 'v', 'i', 'h', '2', 'f', 'w', 't'}
// var idMapper [36][36]byte = [36][36]byte{{'6', 'p', 'e', '5', 'b', 'm', 'o', '0', 'd', '3', 'l', '1', 'r', 'k', 'c', 'q', 'x', 'y', 'g', 'a', 'z', 'u', 'n', '7', 'j', '8', 's', '9', '4', 'v', 'i', 'h', '2', 'f', 'w', 't'}, {'h', 'i', '1', '9', 'x', '7', 'm', 'u', 'j', 'c', 'e', 'd', '4', 'n', 'z', '3', 'a', 't', 'w', '0', 'b', 'r', '8', '5', 'g', 'q', 's', 'y', 'l', 'p', 'v', 'o', '6', '2', 'k', 'f'}, {'e', '3', 'r', '6', 'k', 'z', 'c', '2', 'v', 't', '9', 'n', 'h', '5', 'y', 'w', '8', 'o', 'x', 'b', 'l', '7', '1', 'i', 'p', 'g', 'f', 's', 'q', 'm', '0', 'd', 'u', '4', 'a', 'j'}, {'i', 't', 'r', 's', 'e', '9', 'b', '0', '5', '1', '8', 'g', 'j', 'l', 'a', '6', 'w', '3', 'p', 'k', '4', 'z', 'u', 'o', 'c', 'v', '2', 'n', 'y', 'x', 'q', 'h', '7', 'f', 'd', 'm'}, {'i', 'k', '1', 'q', 'u', '3', 'y', '8', 'b', 'j', 'm', 'f', '7', '2', 't', 'a', 'p', 'r', 'd', 'w', 's', 'o', 'z', 'x', '6', 'h', 'g', 'v', '4', 'c', 'n', '9', 'e', 'l', '5', '0'}, {'8', 'm', 'u', 'b', '5', 's', '9', 'x', 'v', 'i', 'f', 'r', '6', 'k', '2', '4', '3', 'q', '0', 'g', 'h', 'c', 't', 'd', '1', 'a', 'z', 'e', 'o', 'j', 'w', 'p', 'n', 'l', '7', 'y'}, {'t', 'v', 'r', 'k', 'x', 'i', 'u', 'y', '9', 'm', 'c', '5', 'j', '7', 'w', '6', 's', 'z', 'h', '0', 'e', '1', '2', 'g', '3', 'o', 'a', 'l', 'b', '8', '4', 'd', 'f', 'q', 'n', 'p'}, {'y', 'z', 'v', 'b', 'p', 'm', 't', 's', 'n', 'h', 'o', 'w', '8', 'r', 'x', '0', 'j', 'l', 'q', 'k', '9', 'e', 'd', 'a', 'g', '5', '2', '7', 'i', '1', '4', '6', 'u', 'f', '3', 'c'}, {'h', 'l', 'x', 'g', 'q', 'v', 'b', '9', 'p', 'y', '8', 's', 'r', 'j', 'd', 't', '2', '1', 'w', 'a', 'f', 'n', 'o', '7', '3', 'u', '6', 'e', '5', 'i', '0', 'z', 'c', '4', 'm', 'k'}, {'m', '8', 'r', 'f', 'w', 'a', 'n', 'h', 'g', '2', '4', 's', '6', 'v', '1', 'e', 'c', '5', '7', 'u', 'x', '9', 'y', 'k', '0', 'l', 'o', 'q', 'i', 'z', 'd', '3', 'j', 't', 'b', 'p'}, {'a', '8', 'u', 'v', '3', 'e', 'z', 'q', 'c', '1', 'f', '0', 'x', 'g', '2', 'i', '9', 'y', 'p', 'k', 'n', 'o', '7', 'h', 'd', 'w', '6', '4', 'b', 'j', 'r', 's', 't', 'm', 'l', '5'}, {'c', 'u', 'd', 'x', 'y', '9', 'b', 'o', '2', '6', '3', 'j', 'l', 'n', 'i', 'z', 't', '8', '0', 's', '4', 'm', '7', 'f', 'a', '1', 'k', '5', 'q', 'p', 'v', 'r', 'e', 'h', 'w', 'g'}, {'m', 'p', 'x', 'j', 'y', '6', '4', 'a', 'q', 'r', '8', '5', '2', 'b', 'z', 's', '9', 'n', '1', '7', 't', '0', 'h', 'f', 'u', 'd', 'w', 'e', 'o', '3', 'i', 'c', 'g', 'k', 'v', 'l'}, {'b', 'x', 'v', 'p', 'a', 'e', 'd', 's', 'o', 'c', 'r', '6', 'i', 't', 'y', 'm', 'n', '3', '1', 'q', 'j', '5', '8', '9', 'u', 'k', 'w', 'g', 'z', 'f', '2', '7', '0', 'h', '4', 'l'}, {'w', 'x', '8', 'f', '7', 'p', 'g', 'm', '0', 'u', 'h', 'd', '2', 'e', 'o', 's', 't', 'r', 'v', 'q', '9', 'i', 'a', 'z', 'k', 'y', '4', '5', 'c', 'b', 'n', '6', 'j', 'l', '1', '3'}, {'w', 'e', 'x', 'f', 'k', 'z', 'o', 'q', 'l', '5', 'n', 'j', '7', '1', 'm', 'i', '8', 'u', 'y', 'c', '6', 'r', '3', 's', '2', '9', 'v', 'p', 'b', '4', 't', 'g', 'h', 'a', '0', 'd'}, {'g', '8', 'h', 'u', 't', 'q', 'f', 'l', '0', '4', '9', '7', 'o', 'j', 'z', 'i', 'b', 'r', 'w', '5', 'c', 's', 'p', '6', 'n', 'e', '1', 'v', 'x', 'y', 'd', 'k', '3', 'm', '2', 'a'}, {'v', '8', '7', 'b', '9', 'h', 'e', 'u', 't', '1', 'k', '0', '6', 'r', 's', 'n', 'm', 'o', '4', 'q', 'l', 'g', '5', 'p', 'x', 'j', 'i', 'z', '3', '2', 'c', 'a', 'd', 'y', 'w', 'f'}, {'k', 'm', 'b', 'n', 'z', 't', 'x', '0', '9', 'j', '5', 's', 'h', 'r', 'u', 'l', 'p', 'q', 'a', '1', 'i', 'o', '2', 'c', '6', '7', '8', 'e', 'd', '4', 'w', 'f', 'v', 'y', 'g', '3'}, {'o', 'e', 't', 'p', 'k', 'h', 'd', 'w', 'f', '0', 'x', '9', 'j', 'z', 'm', '2', 's', '7', 'u', 'v', 'l', 'n', '3', 'y', '1', '4', '8', 'b', 'c', '6', 'q', 'i', 'a', 'g', '5', 'r'}, {'p', '0', 'n', 'd', 's', '9', 'k', '2', 'w', 'f', '5', 'x', 'z', 'y', 'h', 'b', '4', '7', 'o', 'a', 'u', 'r', 't', '1', 'm', 'j', 'l', 'i', '3', '6', 'g', 'e', 'c', '8', 'v', 'q'}, {'6', 'p', 'y', 'q', '7', 'w', 'l', '0', 'x', '4', 'r', 'k', 'd', '3', 'i', 'h', 'j', 'v', 'a', 't', 'z', 'c', 'b', 's', '1', 'u', '2', 'm', '5', 'e', '9', 'o', '8', 'f', 'g', 'n'}, {'x', 'c', 'j', 'l', 'z', 'n', 'b', 'g', 'd', 'e', '5', 'q', '0', 'y', 'o', 'k', '9', 'u', 'a', 'r', 'v', 'w', 'i', '8', 'f', '2', 'm', 'h', 't', '7', '1', '3', '6', '4', 'p', 's'}, {'m', '6', 'n', 'r', 'e', 'o', 'z', 'a', 'x', 'v', 's', '8', 'd', 'w', '9', '7', '2', '3', '1', 'b', '0', 'i', 'l', 'k', 'q', 'y', 'h', 'f', 'p', '4', 'c', 'u', 'j', 'g', 't', '5'}, {'b', 'c', 'k', 'd', 'l', 'g', 'i', 'u', '0', '5', 'h', '7', 'z', 'n', 't', 'w', '8', 's', '3', 'r', '1', 'p', 'f', 'j', 'a', 'q', '6', '2', '4', 'o', 'm', 'x', 'v', 'e', '9', 'y'}, {'g', 'p', 'f', 't', '6', 'l', 'b', 'd', '1', 'x', 'q', 'm', 'n', 'u', '7', 'a', '4', 'y', 'j', 's', '3', 'v', 'h', 'e', 'i', 'c', '5', '8', 'z', 'o', '2', '0', '9', 'r', 'k', 'w'}, {'h', 'b', 'd', 'p', 'o', '3', '4', 'w', 'l', 'u', 'c', 'v', '9', 'k', '0', '5', 'x', 'j', 's', '6', '7', 'y', 'n', '2', '8', 'e', 'f', 'a', 'g', 'r', '1', 'z', 'i', 't', 'm', 'q'}, {'m', 'n', '9', '5', 'u', 'g', 'z', '3', 'o', '7', 'k', 'f', 'v', '2', 't', 'p', 'q', '4', 'c', 'y', '6', 'w', 'h', 'r', 'x', 'j', 'b', '0', 'a', '8', 's', 'l', 'e', 'd', '1', 'i'}, {'i', 'p', 't', 's', 'u', 'g', 'w', 'z', '1', '3', '0', '2', 'y', '7', 'o', 'x', 'b', '8', '9', 'd', 'r', 'f', 'c', 'a', 'k', 'v', 'l', 'q', 'e', 'n', 'h', 'm', '4', 'j', '5', '6'}, {'w', 'y', '4', '8', 'k', 'x', 'e', 'a', 'b', 'p', 'z', 'i', '6', 'c', 'l', 's', 'v', '3', 'f', 'm', '0', 'n', '9', '2', 'g', '5', 'd', 'h', 'q', 'j', '1', 'r', 't', 'u', '7', 'o'}, {'5', 's', 'c', '7', 'v', '2', 'n', 'z', 'b', 'p', 'h', 'r', 'u', '9', 'l', 't', 'k', 'w', '0', 'a', 'e', 'd', 'y', 'm', 'g', 'i', 'j', 'x', '3', '8', '4', '1', 'o', 'q', 'f', '6'}, {'5', 'q', '3', 'j', 'y', 'u', '1', 'm', 'r', 'k', 'i', 'a', 'l', 'o', 'x', 'c', 'n', '2', 'g', 'b', 't', 'v', 'z', '4', '7', 'd', '8', '9', 'w', 'h', 's', '0', 'e', 'p', 'f', '6'}, {'w', 'k', '7', 'o', 'v', 'b', 't', '0', '6', 'r', 'c', 'z', 'd', 'e', 'l', '9', 'u', 'q', 'f', '8', '5', '2', 'p', 'm', 'x', 's', 'a', 'n', '3', '1', 'i', 'h', 'y', '4', 'j', 'g'}, {'y', 'q', '4', 'f', 'c', '8', '7', 'j', 's', 'o', '2', 'v', 'l', 'g', 'm', 'n', 'b', '5', 'z', 'r', 't', '1', 'a', '9', 'k', '6', 'w', 'x', 'p', 'i', 'h', 'e', '3', '0', 'd', 'u'}, {'h', '2', 'w', 'u', 'n', 'c', 'm', 'o', 'k', 'a', 'r', '6', 'x', 'b', 'y', 'p', 'z', 'i', 'l', '3', 't', 'f', '4', 'v', 'd', 'e', '1', '0', 'g', 's', '9', '8', '7', 'q', '5', 'j'}, {'l', '4', '3', 'v', '9', 'k', 'n', 't', 'i', 'm', 'c', '0', 'f', '8', 'y', 's', 'q', 'e', 'w', 'r', '1', 'j', 'd', 'p', '6', 'h', '2', '5', 'g', 'o', 'a', 'b', '7', 'u', 'z', 'x'}}

// 校验appkey && secret
func CheckAppkeyAndSecret(appkey, appsecret string) bool {

	// 校验appkey是否合法
	checkflag := checkAppKey(appkey)
	if !checkflag {
		return false
	}

	decodeBytes, err := base64.StdEncoding.DecodeString(appsecret)
	if err != nil {
		return false
	}

	desBytes, err2 := decryptDESECB(decodeBytes, []byte("core$#@!"))
	if err2 != nil {
		return false
	}

	req := &message.ServerLicense{}
	if err := proto.Unmarshal(desBytes, req); err != nil {
		return false
	}

	if req.GetAppkey() != appkey {
		return false
	}

	if req.GetNetwork() != getCurrentSerial(config.TOMLConfig.ServerNetwork) {
		return false
	}
	return true

}

// des/ecb/PKCS5Padding解密，同java默认
func decryptDESECB(d, key []byte) ([]byte, error) {

	if len(key) > 8 {
		key = key[:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(d)%bs != 0 {
		return nil, fmt.Errorf("DecryptDES crypto/cipher: input not full blocks")
	}

	out := make([]byte, len(d))
	dst := out
	for len(d) > 0 {
		block.Decrypt(dst, d[:bs])
		d = d[bs:]
		dst = dst[bs:]
	}
	out = pkcs5UnPadding(out)
	return out, nil
}

//明文补码算法
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//明文减码算法
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func _10_to_62(number int64) []int {
	rest := number
	stack := make([]int, 0)
	if rest == 0 {
		stack = append(stack, 0)
	}
	for rest != 0 {
		stack = append(stack, int(rest-rest/36*36))
		rest /= 36
	}
	llen := len(stack)
	num := make([]int, 0)
	for i := range stack {
		num = append(num, stack[llen-i-1])
	}
	return num
}

func getIndexByInt(num []int) int {
	return num[len(num)-1]
}

func getNormalIndexByCharactor(c byte) (int, error) {
	for i := 0; i < 36; i++ {
		if normalMap[i] == c {
			return i, nil
		}
	}
	return 0, errors.New("charactor is unkown")
}

func getIndex(c byte, index int) (int, error) {
	for i := 0; i < 36; i++ {
		if idMapper[index][i] == c {
			return i, nil
		}
	}
	return 0, errors.New("charactor is unkown")
}

func string2Long(appId string) (int64, error) {
	if appId == "" {
		return 0, errors.New("appId is empty")
	}
	index, err := getNormalIndexByCharactor(appId[len(appId)-1])
	if err != nil {
		return 0, err
	}
	var num int64
	for i := 0; i < len(appId)-1; i++ {
		tempInt, err2 := getIndex(appId[i], index)
		if err2 != nil {
			return 0, err
		}
		num = num*36 + int64(tempInt)
	}
	num = num*36 + int64(index)

	num = math.MaxInt64 - num

	if num < 1 || num > math.MaxInt64 {
		return 0, errors.New("num is invalid")
	}
	return num, nil
}

//校验appkey是否合法
func checkAppKey(appId string) bool {
	if appId == "" {
		return false
	}
	index, err := getNormalIndexByCharactor(appId[len(appId)-1])
	if err != nil {
		return false
	}
	var num int64
	for i := 0; i < len(appId)-1; i++ {
		tempInt, err2 := getIndex(appId[i], index)
		if err2 != nil {
			return false
		}
		num = num*36 + int64(tempInt)
	}
	num = num*36 + int64(index)

	num = math.MaxInt64 - num
	if num < 10000000 || num > 99999999 {
		return false
	}
	return true
}

func long2String(appId int64) (string, error) {
	if appId < 0 {
		return "", errors.New("appId is invalid")
	}

	num := _10_to_62(math.MaxInt64 - appId)
	buffer := make([]byte, 0)
	index := getIndexByInt(num)
	end := len(num) - 1
	for i := 0; i < end; i++ {
		buffer = append(buffer, idMapper[index][num[i]])
	}
	buffer = append(buffer, normalMap[index])
	// buffer.append(AppIdMapper.normalMap[index]);
	// return buffer.toString()
	return string(buffer), nil
}

// des/ecb/PKCS5Padding加密，同java默认
func entryptDesECB(data, key []byte) ([]byte, error) {
	if len(key) > 8 {
		key = key[:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	data = pkcs5Padding(data, bs)
	if len(data)%bs != 0 {
		return nil, fmt.Errorf("EntryptDesECB Need a multiple of the blocksize")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func createAppkeyAndSecret(appid int64, secret string, network string) (string, string, error) {

	appkey, err1 := long2String(appid)
	if err1 != nil {
		return "", "", err1
	}

	req := &message.ServerLicense{
		Appkey:    proto.String(appkey),
		Appsecret: proto.String(secret),
		Network:   proto.String(getCurrentSerial(network)),
	}

	buf, err4 := proto.Marshal(req)
	if err4 != nil {
		return "", "", err4
	}

	desEnBytes, err3 := entryptDesECB(buf, []byte("core$#@!"))
	if err3 != nil {
		return "", "", err3
	}

	desEnStr := base64.StdEncoding.EncodeToString(desEnBytes)
	return appkey, desEnStr, err3

	// log.Printf("Debug. parseFrom base64. desEnStr[%v] AppIdentifier[%v]\n", desEnStr, req.String())

}
