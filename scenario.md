# Scenario (Conference Platform)

On the following sections we will be looking at an application developed for conference organizers (stakeholders). We will look at the requirements behind the application and some of the early use cases that conference organizers want to tackle with it. 

The application, as every software, is in constant evolution. Teams will work on adding and delivering new features for conference organizers. 

Let's take a look at a couple of use cases that are vital for conference organizers:
- Home and about page
- Call for Proposals

## Home and about Page

The Home page is quite basic, but it need to show three main pieces of information: 
- Where the conference is happening
- When the conference is happening
- In which phase the conference is (being organized, selling tickets, finished)

Depending on the conference phase, the Home page can show more information such as deadlines for submitting proposals, some highlights from the conference agenda and the conference sponsors.

The About page should cover who the organizers are and how to get in touch with them if you are interested in helping out.

For the initial phase of the conference, when it is being organized, a static site will do, but the moment that we want to ask for potential speakers, we need to provide some functionality for organizers to be able to receive proposals, review them and add the approved ones to the conference agenda. 

## Call for Proposals

When the conference organizers are ready to accept proposals, they should be able to enable the Call for Proposals feature. Because this step happens after defining where and when the confernece will happen, this feature should be implemented first in our Conference application.

The Call for Proposals flow should work like this:
- Conference organizers enable the Call for Proposals flow in the application 
- Potential speakers fill out a form with a session proposal for the conference and submit it for the conference organizers to review
- Conference organizers are in charge of reviewing the submitted proposals and accept or reject them as they see fit. 
- An email is sent to notify the acceptance or rejection of proposals
- Accepted proposals must be added to the conference agenda

This flow requires the conference organizers to have a private page where they can review the submitted proposals. This flow also requires the application to send emails.

When conference organizers are done with selecting proposals (or when the deadline is due), they should be able to close the Call for Proposals feature and publish the conference agenda with all the selected speakers.

## Current state

The team has build a very basic application to host the Home and the About page, and decided that separate services will be created to serve the Call for Proposals functionality and the Conference agenda page.

You can install the current version of the application on a local Kubernetes Cluster by following the [KinD tutorial here](kind/README.md).




