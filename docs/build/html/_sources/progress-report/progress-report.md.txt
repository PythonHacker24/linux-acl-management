# Progress Report for ACL Management GSoC Project

## [9 February 2025] 

The conversation was initiated with Mahmoud Zeydabadinezhad on an Email thread introducing myself and providng compherensive detials about me. This includes all my Social handles relevant to GSoC project, motivation, and a CV for detailed informatiom about my work.   

On the same day, Mahmoud responded with an offer to have a quick call after some discussion about the relevant project. 

I asked about selecting projects given my the organization and Mahmoud suggested me with Project #1 (Securing Linux Storage with ACLs: An Open-Source Web Management Interface for Enhanced Data Protection). 

## [10 February 2025] 

I started researching about the topic and proposed my first designs to Mahmoud. I had primarly 3 questions: 

1. I created this diagram to visualize what I think current researchers use in their network storage. Is this something that you are expecting? 
![first](../_static/first-10-feb.png)
2. As much as I understand, the permissions and access provisions on a high level look like this:
![second](../_static/second-10-feb.png)
Am I right in this case, or are there some additional components to these? 
3. If I consider the above designs to be true, this is the top-level system design I am imagining:
![third](../_static/third-10-feb.png)
Is this something you would consider to be in the right direction? (the primary storage server is actually)

Mahmoud invited me for a Zoom call on Tuesday. We had a successful first meeting and I got introducted to the GSoC organization. 

## [12 February 2025]

Comments from Robert Tweedy arrived on the proposed designs. Here are the comments by Robert Tweedy: 

1. This architecture seems like a reasonable layout; it's not specifically what BMI has (in our case nearly all of our access would be through what's listed as "Remote Access" on that diagram, & we don't have an on-site backup server; we also have some separate storage systems too which aren't directly connected to the main "central" storage & we weren't concerned about these for the purposes of this project as they're not intended for mass-user sharing of data) but is a reasonable architecture to work from and a good testing architecture for the project's application (which shouldn't care too much about the actual network architecture - it's the access to the underlying filesystem with the data that matters, & it can be assumed the server the application's running on will have access to it).

2. This is an ideal way to organize a project directory & looks fine from a purely high-level overview, however it shouldn't be assumed the underlying file organization will match this structure in a real-world scenario where a lab's been around for a long time & might have different/various ways of (dis)organization of their data. The ability of the application to support applying a solution to the following type of request is necessary (as we certainly see this in BMI a lot & presumably other institutions that would benefit from this GSoC project face similar situations in data management): "Hello IT team, this is Dr. Bob from Boblab. I need you to grant Alice (a student in Charleslab) read-only access to this specific directory in my lab space: '/labs/boblab/data_directory_only_boblab_POSIX_group_has_rx_access_to/another_directory/dataForProject/' and read/write access to this different directory: '/labs/boblab/output_directory_only_boblab_POSIX_group_has_w_access_to/another_directory/resultsFromCharleslab/', without changing anyone else's access to these directories & without giving Alice access to directories other than those I listed (and no ability to list the contents of parent directories she has to traverse through to get here); I should also be able to read/write/delete any data Alice uses/creates in these directories because it's my lab space & we can't put it into our institution's shared collaboratory space '/labs/collab/' because <insert_reason_here>. Oh, and Dr. Charles and his Post-doc Danielle should have access to these locations too but have read-write access to everything since they'll be adding data to my 'dataForProject' directory from a collaborator".

3. End-users may not have direct SSH/SFTP access to the underlying storage server nor the ability to establish NFS/SMB mounts depending on security architecture from the IT team & might only be able to interact with the storage through an IT-managed mountpoint on a computational server/node (this is definitely the case in BMI; all access to our /labs volume is done by interacting with the IT-managed mountpoints on the systems, mounted either as NFS or via the BeeGFS client we're using). The ACL Management Application box on the diagram looks fine to me, as long as the "ACL Service" is in fact applying the ACLs to the underlying filesystem & can tolerate changes if made outside of the application directly on the storage servers (likely by the IT team).

## [13 February 2025]

An updated architecture was proposed: 

![updated](../_static/updated-13-feb.png)

**Here is my response with the attached email:**

Clearing out the network architecture would be very helpful in defining the lower-level requirements and designing it. 

I am currently designing the file structure according to the requirements, and so far, I am very happy to say that it's going smoothly. I have come up with a solution but am focusing on optimizing it further to the lowest Linux levels to provide optimal performance. I will let you know soon. I am super excited to take feedback on that one. 

Here are some of my thought processes for your reference:
1. Due to the project's background, More focus must be given to the security portion, and features must be selected carefully. 
2. The project must be highly configurable since it will be used beyond BMI. So, I am thinking of making the configuration as replicable as possible to make it easy for similar architectures to adapt. Since most of this work is done by IT teams, YAML-based configurations make a lot of sense to me.  
3. I would focus on keeping the backend and frontends loosely coupled and developing developer documentation for easily building custom frontends for other users. This loose coupling would allow the development of Android and iOS apps in the future (it's good to be future-proof in this project). As per my experience, developer docs help a lot. The front end that we will be building will be the default one. 
4. The deployment should be simple and highly configurable, making use of tools like Docker to create images of updated versions for different architectures (actually, it might happen that we might use Kubernetes in the future designs, in which case Dockerization would help use a lot).

Apart from this, I want to ask: Can we use Golang for the backend and writing daemons? The selection of Golang is due to its simplicity, performance, concurrency model, powerful, well-maintained standard libraries (by Google themselves), ease of deployment and distribution, etc. It's specifically built for servers; honestly, I can write a whole proposal on why Golang can be a good choice for this project. Although Go has huge accessibility to lower-level features we need, in case we are somewhere limited, I would use C or Rust as per the requirements.

On the same day, a meeting was scheduled where I, Mahmound and Robert discussed about the project in detail. 

## [18 February 2025]

The meeting was held on this day and 3 of us got together in a meeting for the first time and discussed in-depth about the project architecture. I presented slides containing my proposal of solution and sent an email with the updated slides. 

Link to the slides: [Download the PDF](../_static/Introduction-meeting-revised-slides.pdf)

***After this meeting, a lot of emails where exchanged. Each response given by me was based on research done on previous emails by Robert. I considered it to be best if I added email contents here since they contain everything I came upto during the discussion phase. During this time, mentors suggested me not to rush with writing ode and first finalize all the requirements, which I did through these conversations.***

## [22 February 2025]

**Robert Tweedy provided with the follow feedbacks after reviewing the slides:** 

1. Regarding "Proposed System Design for 1st Approach" (slide 6): There's a mention of using Hash Tables for tracking changes to the underlying filesystem to notify the application's backend & frontend processes as part of a responsive change-detection mechanism. I'd highly recommend weighing the benefits of setting this up (& the additional storage required to store said hash tables) vs. just having the application read the current state of the filesystem when it loads the particular branch of the directory tree it's currently working with. Our storage volume is ~3.1PB in size (with many thousands/millions of small files in some directory tree branches) and the separate storage of entire filesystem state in a Hash Table could take up multiple Gigabytes (not to mention the very long initial scan of the filesystem to initialize this hash table). Please let me know if I've misunderstood the intention here.

2. Using Golang for the backend/APIs/daemons/etc. is acceptable, along with Python as needed. Using "C" for lower level operations should be fine as well, but it's also perfectly acceptable to have the application make a call to the getfacl and setfacl commands in the background as needed rather than trying to mess with super low-level filesystem APIs. Javascript for the frontend development is acceptable, and any major JS libraries (jQuery, React.js, Angular, etc.) can be used as well as long as they can be packaged with the application (ie. no calls to external CDNs to load Jquery; it needs to be self-hosted within the application itself).

## [23 February 2025]

**I responded to the above email with the following respose:**

I think proceeding with the second approach would be preferable. It nearly fits with the current implementation and would be more productive. 

In answering the first question, I would need to weigh the comparison between the hashes approach and reading the current state while working on the lower-level system design. I understand the constraints now about our massive amount of storage, and I will address it in the next deck. One of my ideas for this purpose is Event-Driven Architecture, where all the changes would be pushed into the messaging queue for the backend to read and note. Here, RabbitMQ or Apache Kafka can be used. Let me know what you think about this.

For the second point, I will surely keep those things in mind. I will ensure that we avoid the lower levels of the file system API as much as possible. Also, since the front end will be shipped with the application, I will ensure that we don't use CDNs.

In summary, I will develop an updated mid-level combined with a lower-level system design now and send the deck, and then I can start writing the code.

## [26 February 2025]

**Robert had the following comments on my response:** 

Regarding all of the details about the hash tables vs. Event-Driven Architecture, I can see benefit for having EDA set up monitoring the status of requests made by an end user (lets say a user submits a cascading permissions change request that will take hours to complete due to the number of files that need to be changed; this would allow for the application to provide some sort of status/progress update regarding the request until it completes), but I'm not sure that it should be used as an attempt to see all changes taking place on the filesystem in general since these could occur in large quantities for multiple reasons (computational jobs on the research cluster creating new files and immediately adjusting permissions as part of the job's process, for example) and would be very costly to track vs. just reading the current state as-is when a query about a directory/file is made through the application.

If you do feel that including EDA is still worthwhile for the project, then either RabbitMQ or Apache Kafka sound fine to me & you're free to use whichever one you'd prefer to work with.

## [27 February 2025]

**Here is my response to Robert's comments:**

I have realized that I need to do some calculations from my side. The reflection through reading the current state would work fine and be simple enough, but it depends on the volume size. If we have a million files with us, then reading the whole state after a change is challenging and time-consuming. Even if we manage to update the user only with the updates, it would be hard if the user has access to a lot of files and a big permission change has been made. 

I liked the concept of giving a status update to the user if there is a big update from their side. It would definitely improve the user experience and can be done with EDAs. RabbitMQ would be preferable since it's lightweight and simpler in design than the more heavy, advanced capabilities of Kafka, which we would not need.

One mechanism I am currently working on is mirroring the file structure and permissions into a memory-based cache like Redis. Here, all the changes made on the primary file servers would be reflected in the permissions and structures. I am aware of the drawbacks: more storage is required for this, but since it has a very high write and read speed, the user can access the structure faster. (In case the user is reloading repeatedly, it must not put load on the file server.) 

More precisely, when the user triggers a huge permission change affecting 10,000 files, the file server updates it, provides status through EDA, and serves the user. Once the big change is made, the file server updates Redis through notification, and the user dashboard updates with Redis content. Now, the user can log in and reload any amount of time; the file server is not affected. 

I am currently working extensively on these calculations since we will not be able to make many changes in the advanced stages of the project. 

## [1 March 2025]

Robert's response on my last email: 

By "If we have a million files with us, then reading the whole state after a change is challenging and time-consuming" are you referring to reading the state of the entire filesystem, only the directories/files that were updated during that particular transaction, or just the active working directory currently being viewed by the end-user in the application (regardless of whether there were changes in its subdirectories being performed as part of the latest transaction)? I personally think that if a user gets so many files into a single directory (no subdirectories) that it's causing an issue for the underlying filesystem to read the state of that directory then it's a data-management/organization issue on the end-user's side of things & they should interpret the overall slowness as a reason to better organize their data for faster performance (which is out of scope for this GSoC application to rectify).

Additionally, while it is definitely important to make sure that the application makes only reasonable calls to the underlying research filesystem for information when needed so that it doesn't bog it down with unnecessary requests, it shouldn't be left to the application to always assume that it must handle its own filesystem state caching process as opposed to just making requests from the filesystem directly when necessary. In our particular case, our underlying storage system is already making use of dedicated metadata servers so that the actual data storage node is only being accessed when the research data itself needs to be retrieved/written; running commands like ls,chmod,getfacl,etc. are fine & usually do not take that much time to complete for simple non-recursive lookups (though they can take at least a few hours if working recursively to apply changes on a large set of data, like if a researcher decides to recursively make a change to their entire /labs/mylab/ directory, but this is anticipated due to the sheer quantity of changes that would need to be made), and it would be acceptable to us to just have the application perform direct queries from the filesystem using standard APIs or commands like ls and getfacl since we have already tuned our underlying storage infrastructure to handle this at reasonable speeds in most circumstances; the GSoC application doesn't need to be aware of our metadata servers specifically, because any interaction with the filesystem that requires their use is transparently passed along to the metadata servers when needed.

This may be one of the parts of the overall application that you'd consider modularizing, having modules available that could interact with the managed filesystem in different ways (ie. a module that just does direct access via common filesystem commands/APIs as previously mentioned, another module that can perform the type of filesystem tree caching/internal tracking as you've proposed and mentioned you are looking into, and any future modules that might be written by other institutions to better integrate with their specific set-ups) but provides a standardized API that is used by the application so that deployment sites can select the module that best fits their needs. Some sites may have smaller sets of standalone file servers where having the GSoC application store metadata makes sense & would be beneficial to reduce strain on their systems, while others like us already have this resolved in our general storage system structure & would prefer to not need to dedicate additional storage just for the app's own internal cache.

## [3 March 2025]

**My response on Robert's comments:**

Okay, all the above instructions and requirements are clear to me. It seems reasonable, and the Redis idea can be added as a module and developed in later stages. 

As Robert sent in the last emails, I have reviewed the current BMI architecture in detail. The big question is: How much is BMI's current implementation willing to change? 

I know that the fewer the changes, the better the would be for installation on BMI infrastructure, so I want to make sure I work hand in hand with the current implementation. 

Currently, we have an OOD web GUI hosting server and computational job nodes physically connected to storage clusters. We also have LDAP for authentication and management. Since I was told about the current functionality, I can visualize how the configurations may look at a high level. 

I am planning to replace the current OOD server with our backend server, which will serve the GSoC Project frontend to the user, provide computing access/permissions, schedule jobs for the computational nodes, and more. 

We would consider computational nodes and storage clusters hidden and packaged together like a single interface. The storage operations can be performed with Slurm Jobs in the computational nodes, which would be an abstraction. 

Now, we need a mechanism for the access controls that can handle millions of file permissions and 100s of users with fast read and write performance. For this, I am thinking about using LDAP and PostgreSQL in parallel. I want to give a scenario here: 

Think we have Bob as a user. Bob is an intern working with Alice (senior project manager). Alice has access to /lab/alice/biomechanics/confidential/secret.txt, and Bob has his own /lab/bob directory. Suppose Bob is provided access to /lab/alice/biomechanics/confidential/secret.txt. The mechanism I am considering looks like this: First check if Bob has access to /lab: No. Then check if Bob has access to /lab/alice: No. Go on and on; Bob has no access to /lab/alice/biomechanics/confidential. Finally, check for  /lab/alice/biomechanics/confidential/secret.txt, and yes! Bob has the access, and he can work with it. 

Now, this algorithm allows us for nested permissions: if Bob beforehand had access to  /lab/alice/biomechanics, then he must have access to  /lab/alice/biomechanics/confidential/secret.txt and can continue with even reaching check for /lab/alice/biomechanics/confidential/secret.txt 

As per my analysis and comparison with PostgreSQL, LDAP has merits and demerits. It boils down to the performance, so we will use LDAP for higher-level permissions like directories/lab/alice/biomechanics|r (so it has read permission). Hence, if  /lab/alice/biomechanics has 100K files, we don't care; Bob will have access to each of them (all the access requests would be pushed to computational nodes for processing). But for fine-grained permissions like /lab/alice/biomechanics/confidential/secret.txt, we would check PostgreSQL since there might be 10K of such files in them and PostgreSQL is good for this (additionally, Redis can be used here to cache frequently access files for faster performance). 

So now the process becomes simple for LDAP to just check for directories and subdirectories: they would be sure less than the number of files by a huge number and will provide optimal performance here. Also, the updating process becomes simpler. If Bob has access to /lab/alice/biomechanics/confidential/secret.txt and has been told to give access to Aditya, the backend will verify the permission and then update the database. 

LDAP would also be used to check various parameters, from basic stuff like whether the user exists to advanced roles like "Is the user allowed to provide his access to other users?" 

Now, this architecture allows us to have centralized permission and ACL in one place, computational and storage jobs in one place (so we can name them two major packages). Frontend would also be one package (optional for big clients since they have standard APIs, but we would be building a default one) and create optional modules that other developers can build in the community, which would be plug-and-play.

I am not going into the details of the actions that computational nodes would take—they are pretty straightforward and just need designs and approvals. 

Let me know what you think about this. 

After this, I will provide a lower-level system design. Should I start writing the GSoC proposal, where I can propose the updated ideas one by one and adapt them to your requirements and instructions? This will allow me to write the proposal simultaneously, have a standard format for proposing updates rather than writing long emails, and version documents.

**On the same day, I got response from Robert:**

BMI's current cluster infrastructure architecture will not change; what we would do is set up a new server (very likely a VM) for the GSoC application to run on & would expose the front-end via appropriate methods (likely via an existing reverse-proxy system) to make it accessible to end users (which could be a sub-url on an existing domain, so it's important that the application supports running from a reverse-proxy address like "https://cluster.example.edu/datamanager/" or similar, since the base "https://cluster.example.edu/" page may be a landing page with links to various things like the job submission system, a help/documentation system, and now the permissions management system).

Attempting to have this GSoC data permissions management application also handle cluster job submissions is out of scope, as there are already applications for this (Slurm, SGE, PBS/Torque, Moab, etc. for the job queueing, and Open OnDemand for the job front-end interaction via GUI). It's also not desirable for the application to require an underlying system like Slurm/Torque/etc. since there's a possibility this application may be deployed at a smaller scale depending on an end-user's needs (we're planning a central deployment within our department, but a small independent team may have just a single server or two without full job scheduling systems & want to deploy this application directly on them to help manage their local filesystems - this should still be supported by the GSoC application).

If this application is using PostgreSQL/Redis/etc. for its own internal purposes (caching, state management, etc.) then that's fine, however note that when deployed it will be adjacent to and not part of the data-access pipeline for the underlying filesystem; the computational nodes will access the filesystem & its data via standard existing methods, with the node/filesystem validating access permissions using the combination of standard POSIX & extended Linux ACLs applied to the directory path being accessed. They will not be making separate calls out to a PostgreSQL database to do a custom permissions lookup, as that would imply they're accessing the filesystem/data via some sort of middleware layer that controls access to the filesystem & is performing those lookups.

If you want to start writing the GSoC proposal I think that's fine, but I need to defer to Mahmoud to confirm as he has more familiarity with the GSoC program than I do & is more familiar with the overall GSoC process. I would definitely recommend making sure the scope of the project is understood before writing the final proposal, as from our recent emails I worry that there's more potential for scope-creep with this project than had been anticipated & it could inadvertently lead to development time spent focused on things that are actually out of scope of the main goal.

## [4 March 2025]

**My response on Robert's Email: **

Okay, that clears it up: I don't need to jump around computation nodes and consider it part of our product. Our users can have them if they want (and Emory has one), but it's not a requirement in the project. The final question here is, are we going to make use of computational nodes when we deploy it inside the BMI network? I am unsure if I understand the use of computational nodes in your existing infrastructure, but I guess it's for doing heavy computation tasks related to data operations. For example, a permission change of 1 Million files has been requested, so the computational nodes schedule the job, which updates the permissions in the data storage unit. 

However, as mentioned, I would not consider keeping computational nodes as part of the GSoC project. Also, it's clear to me that you need a monolith server (running as a VM) that can handle all operations. This approach works for me, and I will follow it. 

What's going around my mind is using goroutines to do processes that take a long time. 

For the reverse proxy, Nginx would be sufficient. 

The second question is, what does the LDAP server do in the current Emory infrastructure? If I knew its current use, I would be able to adapt it for this GSoC project.

If we use getfacl and setfacl, we will manage permissions on the storage servers themselves, which would create users on the servers. Since we have multiple storage servers, the design for managing permissions across all would differ from what we can do with databases like PostgreSQL. 

To sum everything up till now, if we follow the getfacl and setfacl method, the design would look like this: 
1. The backend server would do the job of serving the front end, exposing the standard interaction APIs, basic authentications like whether the user exists or not, etc., and communicating with the storage nodes. 
2. The storage nodes would run a daemon that exposes follow endpoints: create-users, delete-users, get-acls, set-acls, etc. If a user needs to be created, it would be created on the server itself. Deleting would go the same way. Getting user permissions would use getfacl, and modifying users would use setfacl. We are going to have the owner of the directories, users, and groups in it/ 
3. The daemon responds to the backend server and is served to the front end. 

After clearing this up, things like the process bar with EDA can be discussed. 

If I am still getting out of scope, let me know the portions I should cut out completely. 

**Second Email:**

Currently, I can imagine the LDAP being used just for basic authentication and user sign-in and sign-up. I went through different uses of LDAP (like file permissions), so I was wondering if BMI is using it that way. If it's a basic user auth, then I would be using it as it is in the project. 

## [5 March 2025]

**More questions from me:**

After working through the day, I have circled back on this: How much computational power do we have, or does it even matter for the GSoC project? 

The reason why I stuck around the computational nodes is due to edge cases where millions of files are updated at once. If we consider only the backend server and Linux storage servers to do so, then handling millions of file updates would be resource-intensive and time-consuming. Are we assuming they both are sufficient to carry out the process in the expected time? 

If yes, I would consider a single master-slave architecture (backend as the master and storage nodes as slaves running daemons). All the lookups and updates will go through this pipeline (frontend -> backend -> daemon -> underlying storage).

**Robert's reponse on these questions:**

For the GSoC project the computational power doesn't necessarily matter & it's safe to assume that the application's back-end processes would be running on a sufficiently powered system, and it's good to profile the application to get an understanding of what it spends time on when processing (to try and find ways to improve efficiency) but there's no hard time-limits that we're expecting to be met at the moment as long as the app is responsive to end-user input within a reasonable amount of time (completion of back-end tasks can take as long as needed & shouldn't be blockers on the front-end if at all possible); currently I can dedicate a system with at least 8 cores & 64GB of RAM upon deployment here in BMI (other sites might deploy on lower-specced hardware, but overall this shouldn't be a concern to the GSoC project as it should be generally understood by the end-user that there would be a performance impact depending on the resources available to the application).

The computational nodes are being used for research computational jobs & can frequently be occupied to the point that there are jobs pending in the queue that are waiting for resources to be available; as such, I'd advise against using an architecture for the GSoC application that expects to offload its processing to the nodes (but there's nothing preventing this from being one of the application's many back-end plugins available when the application grows as there may be another research team that is interested in using it like that; just don't make it a focus during this initial development period).

Regarding LDAP: This is being used for basic authentication & user sign-in, along with identifying the Linux/POSIX groups that user accounts are members of, but nothing related to actual file permissions since this is a filesystem-level function. I'd advise not requiring that specific changes be made to the LDAP instance in order for this application to function (ie. no mandatory need to make user groups specific for the app's functions to identify users) and instead leave it up to the IT teams deploying the app to configure things in a way that best suits their environment by providing options in the app to designate one or more (if appropriate) groups to use for things like granting additional admin privileges in the app.

The use of getfacl and setfacl don't involve the creation of users on the servers, they work with the existing users already known to the servers (or UIDs if a user doesn't exist; this type of situation can occur for one reason or another & should be something that the GSoC application can modify/support just in case). Since our servers are connected with our internal LDAP infrastructure, they know the details of all our users as well & there won't be an issue.

## [6 March 2025]

![LDAP](../_static/ldap_6_march.png)

The upper diagram is the high-level system design and below is a transaction scheduler I am working on. 

I took all your emails and compiled each requirement and constraints; and came up with this which fits closest to that list. Let me know what you think. I think this one is much closer to what you are looking for (since I understood the LDAP purpose and ACLs in Linux Servers in better detail). 

The scheduler part is under design right now; but I would love your opinions on this. This design is aimed to improve consistency and resource utilization of the backend server. 

So explain it a bit in here; as soon as the user logs in - they create a session in the backend which is responsible to handle everything the frontend has to do with (this makes our backend server stateful). Each session has its own queue of processes. Processes contain 1 transaction. Here, a transaction is a consistent operation which is true if the operation succeeds and false if it leaves midway (and reverts things back). When a user does some clicks on the frontend and attempts some operations; the frontend shows that they are scheduled (like the upload bar in google drive which uploads in the background and shows the process bar). The session appends the process in a queue. The first in process gets first out and is executed by the transaction executer. If a user quits midway (due to any reason; like an emergency power off on desktop); the backend keeps executing the process as it's scheduled. It executes all of it and when the queue is empty; exits and the session is successfully closed. 

Moving more advanced; if we have a local database that keeps a copy of the queue; it would provide more redundancy in case the backend application crashes midway or server goes offline. This will allow us to keep track of the queue and not end up messing up the permissions in these edge cases. I know that it will not be able to recover 100% but this small integration will improve the consistency drastically.

## [7 March 2025]

**Robert's comments on latest architecture:**

This latest diagram looks like it's nearly there, but I do have the following comments about it:

* The backend server shouldn't be modifying the LDAP server (unless this is an optional plugin for future development, used by smaller-scale deployments that may not already have an existing LDAP infrastructure); it would only be binding to/reading from it so that it can recognize network users & groups (& other details that may be present in LDAP, as appropriate) but shouldn't be making its own changes.

* Your design may already handle this, but I wanted to make sure this was accounted for just in case - The daemon/API running on the file storage server should be capable of binding to a Unix socket too rather than explicitly requiring a TCP/IP socket (or, it should at least tolerate a localhost-only connection); there's a possibility that some sites may want to deploy the GSoC backend on the same system as the storage server depending on the resources they have available & wouldn't want to open multiple firewall ports (one for the backend server's daemon & another one for the storage server daemon).

* Regarding the example you fleshed out in the email: it might be a good idea to make the "revert changes if a transaction fails" functionality an option that the end-user can toggle (on, off, or on up to a certain limit of files/tasks & off above that limit) rather than always being enabled; some users may be fine with leaving a partial transaction completed as-is (since technically they could use the GSoC app to fix any issues/create a new request to complete the rest of the original transaction after fixing the underlying issue) due to the amount of time it would take for the changes to be reverted. There's also the aspect of having to keep track of the original & changed state for each file while a transaction's underway, which could be very memory/storage intensive for large transactions.

Other than those few observations, everything else you've mentioned in this latest email sounds good; the session scheduling mechanism definitely sounds like a good thing to have since it will allow users to not need to leave their systems actively open & watching a large transaction complete, and I'm glad you're excited about working on that aspect of it since it'll be of significant value to the overall GSoC application.

**On the same day, I proposed the following idea:**

Greetings of the day, 

Okay, now I am only short of one thing here: understanding the LDAP server's purpose. I understand that it stores basic user and group data that Emory already has and is used for basic authentication and checks. 

My understanding is that the setfacl and getfacl commands would use the LDAP server to manage users and groups since they are not created on the Linux File Storage Servers. Please correct me if I am wrong. 

Are we allowing the creation of new users or new groups from the front end?

I am imagining a scenario here. A new user can visit https://gsoc-app.emory.edu/signup and create a new account. This info will go on the LDAP server, and the new user will be allocated no privileges to access any files. This user now goes to work with Prof. X. Prof. X has access to lab/MRI-scan data, which he wants the new user to access. He can then give permission to the new user from his logged-in portal. 

This allows new users to create an account on the app with no privileges; others with their own files and permissions can assign the user their permission. So there is no security issue or unauthorized access due to the least privilege but the ability to create new user accounts. As far as I can understand, this information would be stored on the LDAP server. 

Also, users would be allowed to change their passwords, delete accounts, forget passwords, etc., on their own, and there would be no need to contact the IT team for these operations. 

If we do not allow user creation, we are assuming that we are not creating new accounts and are only managing existing users listed on the LDAP server. In this case, the IT team would be responsible for adding new users, changing account passwords, deleting users, etc. 

Let me know about this, as it will clear up the authentication part, and we will almost reach a consensus. 

Regarding daemon connectivity, it sounds good, and I will design the daemon to handle Unix Sockets. I agree that a few sites would consider deploying everything in the same server (Linux File Storage Server) and I would make sure we make provisions for that. In this case, I would call it a monolith deployment where a single server does everything - frontend handling, backend operations, authentication, and storage services. When a site has an infrastructure pictured in the last email's diagram, it can be called an orchestrated deployment where multiple Linux Storage Servers are orchestrated simultaneously. We can provide docs for both kinds of deployment. 

Returning to the process queue and transaction mechanism, I am glad you liked it. I completely agree that the transaction rollback feature can be kept optional since it would require an additional database running with the backend server. Three options seem fine: ON (where when a transaction takes place, during the execution time, it's logged and mirrored in the database no matter how big it is), OFF (where there is no transaction process state saved and would be prone to partial transactions if something wrong happens) and Partial (where a limit would be set on how big a transaction can be eligible for rollback and beyond that if a user initiates a transaction, they would be prompted with the warning of no rollback before confirming). 

This log can be either deleted after the end of the transaction or stored in the database if configured to be. In the prior case, this will allow us to revert back if the transactions fail in the execution period and be true only if the full operation is completed. In the later ones, a transaction can be reverted back anytime (like git). Again, limits can be implemented here (in case we need to skip big transactions). 

Let me know what you think. After the LDAP part is clear, I would like to finalize the design and then work on the lower levels of software implementation, which I would use in the GSoC proposal to sum it up.

## [8 March 2025]

**Robert's comments on my reponse:**

I highly recommend that you check the Manual pages for the setfacl & getfacl commands to get an understanding of what they do, as they only manage extended file permissions & aren't user-management commands; an online version of them is available here, but it should also be available from any Linux system with the appropriate acl package installed via the man command (the package name differs depending on distribution):
https://linux.die.net/man/1/setfacl
https://linux.die.net/man/1/getfacl

As far as we're concerned in our department here at Emory, we will not allow user creation through this GSoC application; other deployment sites might have a use for this so it could certainly be a plugin for the application, but I'd rank this as low priority for now & shouldn't be focused on if there's other tasks to still complete for the application's main development.

Our LDAP server is just used for centralized management of user accounts instead of using local-only files like /etc/[passwd|shadow|group]; worrying about its specifics is out of scope for the GSoC project. All the GSoC project needs to be concerned with is that:
LDAP exists & that it's able to query the LDAP server given the appropriate info provided by the administrator setting up the application (LDAP server address(es), bind user/dn, bind password, base search DN, object filter, attribute mapping, etc.) to validate whether a user account exists & is authorized to log in to the application.
The underlying servers (backend/frontend/storage servers/other servers on the network that are outside of the scope of this project) are already configured to use the central LDAP servers to look-up user names/account UIDs/group membership/etc. and match them together for system-level purposes, so the GSoC application does not need to worry about this. Put another way, all the accounts in LDAP "exist" on all the servers, so a command like id aditya or getent passwd account_uid could be run on any server & would return the same result regardless of the server it was executed on, even though the account doesn't exist in the server's /etc/passwd file.

There should be pre-existing open-source projects that can provide the basic features of this functionality within the GSoC application, so I'd recommend just reviewing those & selecting the one that you feel is the best fit for inclusion with the application (make sure to review their OSS licenses & ensure that they're compatible with the ones required as part of the GSoC!) rather than trying to code an LDAP client implementation from scratch.

## [9 March 2025]

**My updated design:**

I went through the getfacl and setfacl docs, went around and tested it on my Arch Linux system, and read about it in depth. 

Based upon that, I have modified my system design as shown below:

![ldap_updated](../_static/updated_ldap-9-march.png)

There are no provisions for updating anything on the LDAP server for the backend (unless we build a plugin for the GSoC application for this purpose, which I am leaving as a low-priority task, as mentioned). It's clear to me now that the LDAP server should be used in BMI's network, and I will make sure the GSoC application doesn't disturb it. This will also be explained in detail in the documents. 

Let me know if there are any more improvements; I will add them, too. After this, I would be working on designing the lower levels and testing them out by starting to write code, so let me know about this, too. 

Also, since our backend is written in Golang, we can use https://github.com/go-ldap/ldap (MIT Licensed) for LDAP integration. I will also be testing it out soon. 

In addition, I have been working on deciding the interaction protocol between the daemon and the backend server. gRPC catches my attention for remote connections since our endpoints will be well-defined. As mentioned, Unix domain sockets seem to be ideal for local connections.

***After this point, I started working on my GSoC Proposal and interacted with my mentors for feedbacks and updated it as per their recommendations.***

## [21 March 2025]

**I continued with my questions which came during proposal writing and prototyping:**

Sure, no worries. I am writing a new chapter in the proposal for the prototype, which I am working on. The good news is that the proposed session management algorithm is working well, and we can proceed. GSoC proposal submission starts tomorrow and ends on 8 April, so I would like to work more on the proposal and the prototype and get your feedback before submitting it. 

Also, I have one more question: We have multiple Linux File Storage Servers in orchestrated deployments. All those different servers have filesystems where users have permission and access assigned. Should the portal show the user that all the files are in the same server space and hide the orchestration (like virtual machines), or should I provide options about which server to store the file in?

In the prior case, when users open the portal and log in, they would see the files and directories listed (no matter which server they are from). When they do to /test-result and say it was in server-3, the backend would automatically move to the server-3 /test-result content and list it. The session-management module would store all this information. 

In the latter case, the user would select a server like "server-1", "server-2", "server-3", etc., and then the backend would only list those files and work with that server. 

This is important to me since I am now starting to build the prototype of the transaction executor and will be connecting it to the session management module I have built. 

I am looking forward to progressing on the prototype while finalizing the proposal. Even after the submission, I can keep working on it since the modules in the final project will take some of these prototype components and build upon them. It would speed up the final project and allow me to work on more features.

**Robert's response:**

In the case where there's multiple different servers in an orchestrated deployment serving their own unique filesystem, it's possible that all of their filesystems may be mounted on the backend server such that they're accessible simply by navigating to a specific directory (ie. /mnt/server-1, /mnt/server-2, etc.) and I would hope that the end-user administrator deploying the software would configure it to match how their users expect to see the filesystems on their network. As such, I don't think any special options need to be added to the application to explicitly choose a server, as this should be implicit based on the path of the file that the user's navigating to.

I do see where this may be of use for a site that's aiming to have a single deployment of this software which supports multiple different departments at once as opposed to each individual department having their own deployment, but I wouldn't worry about that too much for this initial phase & it can always be a future improvement that's made if needed.
