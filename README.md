# Protoflow

> Put it in a box and forget about it

The fastest way to build scalable, durable, and testable microservice.

Note: This project is still in **pre-alpha**. We are developing this project in the open to get feedback from the community. Please feel free to open an issue or PR if you have any feedback.

## Install
```shell
curl https://raw.githubusercontent.com/protoflow-labs/protoflow/main/scripts/install.sh | sh
```
or
```shell
go install github.com/protoflow-labs/protoflow@latest
```

## Run
Visually build your microservice with Protoflow Studio:
```shell
protoflow studio
```

## Demo
[![youtube demo screenshot](http://img.youtube.com/vi/ZnUyUbh-Xp8/0.jpg)](https://www.youtube.com/watch?v=ZnUyUbh-Xp8)

## Hack
### Backend
```shell
go run main.go studio --dev
```

Backend will proxy to the frontend, so you can just open http://localhost:8080 to see the frontend.

### Studio
```shell
npm install
npm run dev
```

Starts the studio frontend at http://localhost:8000.

### Haphazard notes
- "a framework to easily build small/modular grpc services" -> yes, i want to make it so you bring your code:  
  ``` js
      function foo(x) {
        y = process(x);
        store(y);
      }
      function bar() {
        foo(1);
      }
  ```
  and then with next to no refactoring, foo becomes a scalable, durable function that can be used elsewhere as a buildable lego block somewhere else  
- zapier has the business ops market pretty well captured, but as like a general developer, there is a lot of code that could be put into a box and forgotten about  
- ie. you define the api contract for the block, you fulfill it with an implementation, and then it just exists and keeps working until something breaks, and then you have all eyes on that very clearly defined block  
- where it starts to get really interesting is when you start storing IO to these blocks and that production IO becomes your test cases for the block functioning as expected. if the output starts changing for a given input, then you know immediately you have a breaking changing (and the semver for the block needs to be adjusted accordingly)  
- the way that I am thinking about it is that it is an attempt at making junior developers think more functionally while programming, and the functional way of programming is the easier way to program because you have given them enough lego blocks to build with  
- you make global state less appealing because you can connect your postgresql resource block to your function block so you can make database queries  
- you make the blast radius of code smaller when you isolate logical components, it makes code analysis more feasible (than the impossible task that it is right now). audits are performed for frequently, but less context is needed (since you know what inputs and outputs exist for a given block)  
