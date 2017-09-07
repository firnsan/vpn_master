package main
import (
    "crypto/md5" 
    "encoding/hex" 
	"time"
	"math/rand"
)


const (
    KC_RAND_KIND_NUM   = 0  // 纯数字
    KC_RAND_KIND_LOWER = 1  // 小写字母
    KC_RAND_KIND_UPPER = 2  // 大写字母
    KC_RAND_KIND_ALL   = 3  // 数字、大小写字母
)
 
// 产生随机字符串
func randString (size int, kind int) string {
	metadata := [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}
    result := make([]byte, size)
	kindForEach := kind

    is_all := kind > 2 || kind < 0
    rand.Seed(time.Now().UnixNano())
    for i :=0; i < size; i++ {
        if is_all { // random ikind
            kindForEach = rand.Intn(3)
        }
        scope, base := metadata[kindForEach][0], metadata[kindForEach][1]
        result[i] = uint8(base+rand.Intn(scope))
    }
    return string(result[:])
}


type User struct {
	Id int `json:"-"`
	Name string `json:"name"`
	Password string `json:"password"`
	PasswordEncoded string `json:"-"`
	Zone string 
	CreateTime time.Time 
	ExpireTime time.Time 
}

func NewRandUser() *User {
	name := "user"+randString(4, KC_RAND_KIND_UPPER)
	password := randString(8, KC_RAND_KIND_ALL)
	
	w := md5.New();
	w.Write([]byte(password)) 
	passwordEncoded := hex.EncodeToString(w.Sum(nil))

	return &User{Name:name, Password:password, PasswordEncoded:passwordEncoded }
}

