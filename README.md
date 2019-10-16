## CScoupler

A platform for companies and students to come in touch with one another, easing the job/internship hunt for students and easing the quest of finding qualified students for companies.

##### Endpoints:
* /signin -> 
    + POST request -> 
        + expects json -> 
            + used to sign in and get a jwt
---
* /signup/company -> 
    + POST request -> 
        + expects json -> 
            + used to create a new company and first representative account
---
* /companies/{company_id} -> 
    + GET request -> 
        + used to fetch a company
--- 
* /signup/representatives/invite/{company_id}/{invitelink_id} -> 
    + POST request -> 
        + used to create a new representative. This page is reached by getting invited by an invitelink.
--- 
* /representatives/{representative_id} -> 
    + GET request -> 
        + used to fetch a representative
--- 
* /representatives/invitelink/ -> 
    + POST request -> 
        + used to create a new invitelink. This also returns the created invitelink as response.
--- 
* /representatives/projects/ -> 
    + POST request -> 
        + used to add a new project to the company the representative represents.
---
* /signup/student -> 
    + POST request ->
        + expects form-data with key=resume,val=[a pdf file] and key=studentData, val=[json body containing student data]
            + used to create a new student
---
* /students/{student_id} ->
    + GET request ->
        + used to fetch a student 
---
* /students/edit/{student_id} ->
    + PUT request ->
        + expects form-data with key=resume,val=[a pdf file] and key=studentData, val=[json body containing student data]
            + used to update a student
---