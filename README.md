# วิธี Compile

1.  Clone git ลงมาที่่เครื่อง
2.  รันคำสั่ง `go mod download`
3.  รันคำสั่งเพื่อ compile เป็น binary `go build -o server .`
4.  ก้อปปี้ไฟล์ใน folder `server` ไปรันบนเครื่อง server ได้เลย

# หรือถ้าขี้เกียจก็เอา Dockerfile ไป build ใช้งานได้เลยครับ