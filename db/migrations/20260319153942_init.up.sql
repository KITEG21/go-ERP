CREATE TABLE departments (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE workers (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  department_id INT NULL,
  salary REAL NOT NULL,
  hire_date DATE NOT NULL,
  CONSTRAINT fk_worker_department FOREIGN KEY (department_id) REFERENCES departments(id)
);

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  password TEXT NOT NULL
);

CREATE TABLE attendance (
  id SERIAL PRIMARY KEY,
  worker_id INT NOT NULL,
  check_in TIME NOT NULL,
  check_out TIME NOT NULL,
  date DATE NOT NULL,
  status INT NOT NULL,
  CONSTRAINT fk_attendance_worker FOREIGN KEY (worker_id) REFERENCES workers(id)
);

CREATE TABLE payrolls (
  id SERIAL PRIMARY KEY,
  worker_id INT NOT NULL,
  month TEXT NOT NULL,
  base_salary REAL NOT NULL,
  bonus REAL NOT NULL,
  deductions REAL NOT NULL,
  net_salary REAL NOT NULL,
  status INT NOT NULL,
  CONSTRAINT fk_payroll_worker FOREIGN KEY (worker_id) REFERENCES workers(id)
);