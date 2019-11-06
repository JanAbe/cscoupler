CREATE TABLE IF NOT EXISTS "User" (
    user_id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    "role" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "Company" (
    company_id UUID PRIMARY KEY,
    "name" TEXT NOT NULL,
    information TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "Student" (
    student_id UUID PRIMARY KEY,
    university TEXT,
    skills TEXT[],
    experiences TEXT[],
    short_experiences TEXT[],
    wishes TEXT,
    "status" TEXT NOT NULL,
    "resume" TEXT,
    ref_user UUID REFERENCES "User" (user_id)
);

CREATE TABLE IF NOT EXISTS "Representative" (
    representative_id UUID PRIMARY KEY,
    job_title TEXT NOT NULL,
    ref_user UUID REFERENCES "User" (user_id),
    ref_company UUID REFERENCES Company (company_id)
);

CREATE TABLE IF NOT EXISTS "Project" (
    project_id UUID PRIMARY KEY,
    "description" TEXT NOT NULL,
    compensation TEXT NOT NULL,
    duration TEXT NOT NULL,
    recommendations TEXT[],
    ref_company UUID REFERENCES Company (company_id)
);

CREATE TABLE IF NOT EXISTS "Message" (
    message_id UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now(),
    sender UUID REFERENCES "User" (user_id),
    receiver UUID REFERENCES "User" (user_id),
    body TEXT NOT NULL,
    ref_project UUID REFERENCES Project (project_id)
);

CREATE TABLE IF NOT EXISTS "Address" (
    address_id UUID PRIMARY KEY,
    street TEXT NOT NULL,
    zipcode char(7) NOT NULL,
    city TEXT NOT NULL,
    "number" TEXT NOT NULL,
    ref_company UUID REFERENCES Company (company_id)
);

CREATE TABLE IF NOT EXISTS "Invite_Link" (
    invite_link_id UUID PRIMARY KEY,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    expiry_date TIMESTAMP,
    used BOOLEAN NOT NULL,
    ref_representative UUID REFERENCES Representative (representative_id)
);