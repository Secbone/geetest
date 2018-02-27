package geetest

import (
    // "fmt"
    "time"
    "strconv"
    "net/http"
    "net/url"
    "encoding/json"
    "crypto/md5"
    "encoding/hex"
    "math/rand"
    // "io/ioutil"
)

type RegisterResult struct {
    Success int
    Challenge string
    Gt string
    NewCaptcha bool
}

type Tester interface {
    Register() RegisterResult
    Validate(bool, string, string, string) bool
}

type tester struct {
    Id          string
    Key         string
    ApiAddr     string
    RegisterPath  string
    ValidatePath    string
    Request     *http.Client
}

type RegisterResp struct {
    Challenge string `json:"challenge"`
}

func (t tester) Register() RegisterResult {
    addr := t.ApiAddr + t.RegisterPath + "?gt=" + t.Id + "&json_format=1&sdk=Go_0.0.1&client_type=unknown&ip_address=unknown"

    resp, err := t.Request.Get(addr)

    var (
        challenge string
        success int
    )

    if err != nil {
        // fallback
        success = 0
        challenge = t.fallbackChallenge()
    } else {
        data := RegisterResp{}
        json.NewDecoder(resp.Body).Decode(&data)

        success = 1
        challenge = Md5(data.Challenge + t.Key)
    }

    return RegisterResult{
        Success: success,
        Challenge: challenge,
        Gt: t.Id,
        NewCaptcha: true,
    }
}

type ValidateResp struct {
    SecCode string `json:"seccode"`
}

func (t tester) Validate(fallback bool, challenge string, validate string, seccode string) bool {
    if fallback {
        // fallback
        return Md5(challenge) == validate
    } else if Md5(t.Key + "geetest" + challenge) != validate {
        return false
    } else {
        data := url.Values{}
        data.Add("gt", t.Id)
        data.Add("seccode", seccode)
        data.Add("json_format", "1")

        resp, err := t.Request.PostForm(t.ApiAddr + t.ValidatePath, data)

        if err != nil {
            return false
        }

        result := ValidateResp{}
        json.NewDecoder(resp.Body).Decode(&result)

        return Md5(seccode) == result.SecCode
    }

}

func (t tester) fallbackChallenge() string {
    rand.Seed(time.Now().UnixNano())
    rnd1 := Md5(strconv.Itoa(rand.Intn(100)))
    rnd2 := Md5(strconv.Itoa(rand.Intn(100)))

    return rnd1 + rnd2[0:2]
}

func New(id string, key string) Tester {
    client := &http.Client{}

    return tester{
        Id: id,
        Key: key,
        ApiAddr: "http://api.geetest.com",
        RegisterPath: "/register.php",
        ValidatePath: "/validate.php",
        Request: client,
    }
}

func Md5(source string) string {
    hasher := md5.New()
    hasher.Write([]byte(source))
    return hex.EncodeToString(hasher.Sum(nil))
}
