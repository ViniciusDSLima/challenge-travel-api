INSERT INTO users (name, email, password, role, created_at, is_active)
VALUES (
           'Administrador',
           'admin@example.com',
           '$2a$10$PXQhIHyYoOtZskWmf2mTnugAh.u7DEQDWJshLXuMcy9TDTvrBEWuy', -- senha: admin123
           'ADMIN',
           NOW(),
           true
       );

INSERT INTO users (name, email, password, role, created_at, is_active)
VALUES (
           'Usuário Teste',
           'user@example.com',
           '$2a$10$YujpaqSWr4q4RF8HxSl/1.gyRtl/FqSiEzUQjWOleJBgMVFgU41lq', -- senha: user123
           'USER',
           NOW(),
           true
       );

