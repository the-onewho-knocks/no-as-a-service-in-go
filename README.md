# ❌ No-as-a-Service in GO 
<img width="1332" height="526" alt="Screenshot 2026-01-17 235740" src="https://github.com/user-attachments/assets/60ad57b7-d060-4b0b-8100-ac06b41bb6ad" />

Ever needed a graceful way to say “no”?
This tiny API returns random, generic, creative, and sometimes hilarious rejection reasons — perfectly suited for any scenario: personal, professional, student life, dev life, or just because.

Built for humans, excuses, and humor.

---

## 🚀 API Usage

Method: GET
Rate Limit: None

🔄 Example Request
```
GET /no
```
✅ Example Response
```json
{
  "reason": "I have to keep the couch from floating away, it's an important job."
}
```
---

## 📁 Project Structure

```
no-as-a-service-in-go/
│
├── main.go            # Entry point, HTTP server, routing
├── go.mod             # Go module definition
├── go.sum             # Dependency checksums
│
├── no.json       # List of rejection reasons (API data)
│
├── .env               # Environment variables (PORT, etc.)
│
└── README.md          # Project description and usage
```
---

## 👤 Author

Created with creative stubbornness by [hotheadhacker](https://github.com/hotheadhacker)

Ported to GO by [Hardik](https://github.com/the-onewho-knocks?tab=repositories)


---

## 📄 License
MIT - as original project.

