## CScoupler

A platform for companies and students to come in touch with one another, easing the internship hunt for students and easing the quest of finding qualified students for companies.

---

#### Architecture

##### Docker
To run the application locally, follow the following steps:
1. Download + install git
2. Download + install docker
3. Download + install docker-compose
4. git clone this repository (git@github.com:JanAbe/cscoupler.git)
5. cd into the newly created repository
6. Open a terminal and type ```docker-compose up```
7. Navigate to http://localhost:8080

##### Frontend
The frontend has been made using Vue.js and various tools
such as vue-cli, vue-router and tailwind.css

This client consumes all the endpoints that the backend app exposes
and acts as a userfriendly way to interact with the system.

The webapp is viewable on both large, medium and small screens and is (hopefully) intuitive to use.

##### Backend
The backend code-base follows the hexagonal architecture.

The handlers package corresponds to the adapters part of the hexagonal architecture.
In this package code transforms incoming data / requests to a format that is understood by the services package.

The services package corresponds to the application part of the hexagonal architecture.
This package contains all the use cases that the application supports + some helper methods.

The domain package corresponds to the domain part of the hexagonal architecture.
In this package all the domain objects are defined aswell as the business rules.

##### Database
PostgreSQL is used to persist/store all data.

---

##### Endpoints:
* /signin -> 
    + POST request -> 
        + expects json -> 
            + used to sign in and get a jwt
<br>
* /signup/company -> 
    + POST request -> 
        + expects json -> 
            + used to create a new company and first representative account
<br>
* /companies/{company_id} -> 
    + GET request -> 
        + used to fetch a company
<br>
* /signup/representatives/invite/{company_id}/{invitelink_id} -> 
    + POST request -> 
        + used to create a new representative. This page is reached by getting invited by an invitelink.
<br>
* /representatives/{representative_id} -> 
    + GET request -> 
        + used to fetch a representative
<br>
* /representatives/invitelink/ -> 
    + POST request -> 
        + used to create a new invitelink. This also returns the created invitelink as response.
<br>
* /representatives/projects/ -> 
    + POST request -> 
        + used to add a new project to the company the representative represents.
<br>
* /signup/student -> 
    + POST request ->
        + expects form-data with key=resume,val=[a pdf file] and key=studentData, val=[json body containing student data]
            + used to create a new student
<br>
* /students/{student_id} ->
    + GET request ->
        + used to fetch a student 
<br>
* /students/edit/{student_id} ->
    + PUT request ->
        + expects form-data with key=resume,val=[a pdf file] and key=studentData, val=[json body containing student data]
            + used to update a student