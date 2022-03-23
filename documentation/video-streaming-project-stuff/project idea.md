# Project Idea 

A fully featured _**on-demand**_ video lecture streaming service

In very simple words, here are the key features of our project

- A professor creates their account
	- After which they can create classrooms based on the subjects/classes they take
	- Each classroom will have a unique code

- A student creates their account
	- After which they can _join_ a particular class with the help of the unique code

- A professor can upload recordings of their lecture

- The video is then processed by the server, replicated in 2-3 qualities and segmented into smaller `.ts` chunks along with a manifest or an index file knows as `.m3u8` file and make them available for **on-demand streaming**.

- The student logs in to their account, views the 