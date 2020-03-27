const jwt = require("jsonwebtoken");

module.exports = (req, res, next) => {
  try {
    // Send token in req in format: Bearer <token>
    const token = req.headers.authorization.split(" ")[1];
    req.body.user = jwt.verify(token, "secret_this_should_be_longer");

    const expiresIn = (req.body.user.expiresIn);
    const iat = (req.body.user.iat);

    if(isNaN(expiresIn) || isNaN(iat)){
      throw new Error("expiresIn or iat is NaN.");
    }

    // Time validation
    const now = Date.now() / 1000;
    if ((now - iat) > expiresIn) {
      throw new Error("Time Expired!");
    }

    console.log("AT AUTH:", req.body.user);
    next();
  } catch (error) {
    console.log(error);
    res.status(401).json({ message: "Auth failed!" });
  }
};
