# Buat migrasi
muhammadrizkisetyanto@MacBook-Air-Muhammad arabiya-syari-fiber-1 % migrate create -ext sql -dir internals/database/migrations create-table-user-profile
/Users/muhammadrizkisetyanto/Documents/arabiya-syari-fiber-1/internals/database/migrations/20250226092703_create-table-user-profile.up.sql
/Users/muhammadrizkisetyanto/Documents/arabiya-syari-fiber-1/internals/database/migrations/20250226092703_create-table-user-profile.down.sql


# Up-Down migrasi
Masukan number untuk mengarahkan mau berapa banyak file yang di up/down. Misal 1,2,3 dst
**UP**
migrate -database "postgresql://postgres.jvpemvuvqnmzasjfihys:Wedangjahe217312@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres" -path internals/database/migrations up

**DOWN**
migrate -database "postgresql://postgres.jvpemvuvqnmzasjfihys:Wedangjahe217312@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres" -path internals/database/migrations down

# Dirty migrasi
muhammadrizkisetyanto@MacBook-Air-Muhammad arabiya-syari-fiber-1 % migrate -database "postgresql://postgres:qXdMRsMSGEgQvVrLuBjmUAGkytJwsaWk@trolley.proxy.rlwy.net:59123/railway" -path internals/database/migrations force 20250221005048

# Masuk database
muhammadrizkisetyanto@MacBook-Air-Muhammad arabiya-syari-fiber-1 % PGPASSWORD="qXdMRsMSGEgQvVrLuBjmUAGkytJwsaWk" psql -h trolley.proxy.rlwy.net -p 59123 -U postgres -d railway


# Refresh port
kill -9 $(lsof -t -i:8080)
kill -9 $(lsof -t -i:8080) && go run main.go


# Hapus Versi Migrasi yang Bermasalah dari Database
Jika ingin menghapus versi 20250306232632 dari database secara manual, jalankan perintah SQL berikut di PostgreSQL:

DELETE FROM schema_migrations WHERE version = 20250306232632;

Kemudian jalankan ulang migrasi:


# JWT
muhammadrizkisetyanto@MacBook-Air-Muhammad arabiya-syari-fiber-1 % export JWT_SECRET=rahasia_dong

muhammadrizkisetyanto@MacBook-Air-Muhammad arabiya-syari-fiber-1 % echo $JWT_SECRET

rahasia_dong


# Mencari kata 
muhammadrizkisetyanto@MacBook-Air-Muhammad arabiya-syari-fiber-1 % grep -r "subategories_id" .


# Midtrans 
https://simulator.sandbox.midtrans.com/


# Password
Wedangjahe217312!