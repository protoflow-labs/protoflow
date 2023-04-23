import { register } from "./register.js";
import express from "express";

const app = express();

register().then((registry) => {
  app.get("/:func", async (req, res) => {
    const func = registry[req.params.func];
    if (func) {
      const result = await func(req, res);
      res.status(200).send(result);
    } else {
      res.status(404).send("Not found");
    }
  });

  app.listen(3080, () => {
    console.log("Listening on port 3000");
  });
});
