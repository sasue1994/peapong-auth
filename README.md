# peapong-auth

#Week 1 ลองเขียน อ่าน request response + เชื่อมฐานข้อมูลอย่างง่าย


#Week 2 เป้าหมาย ปรับโครงสร้างให้เป็นมาตราฐาน

peapong-auth-service/
├── cmd/
│   └── api/
│       └── main.go         <-- จุด Main
├── internal/               
│   ├── app/                <-- App Context ( Dependencies และ Routes)
│   └── auth/               
│       ├── api/            <-- Handler (Controller)
│       ├── service/        <-- Bussiness Logic
│       └── repository/     <-- ต่อ DB
├── pkg/                    <-- Utility
├── go.mod