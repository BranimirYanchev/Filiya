import React, { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import styles from "./SignUp.module.scss";
import Button from "../../components/Button/Button";
import { useAuth } from "../../contexts/AuthContext";
import { GOOGLE_LOGIN_URL } from "../../utils/api";

export default function SignUp() {
  const navigate = useNavigate();
  const { register, isAuthenticated, user } = useAuth();
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [repeatPassword, setRepeatPassword] = useState("");
  const [signUpError, setSignUpError] = useState("");

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/feed');
    }
  }, []);
  useEffect(() => {
    if (isAuthenticated) {
      console.log(user);
      navigate('/feed');
    }
  }, [isAuthenticated, navigate]);

  const signUp = async (e) => {
    e.preventDefault();
    setSignUpError("");
    
    try {
      await register({
        fullName,
        email,
        password,
        repeatPassword,
      });
      
      navigate('/feed');
    } catch (error) {
      const errorMessage = error.message || error;
      if (errorMessage.includes("unprovided fields")) {
        setSignUpError("Please fill in all fields");
      } else if (errorMessage.includes("invalid email")) {
        setSignUpError("Please enter a valid email address");
      } else if (errorMessage.includes("invalid full name")) {
        setSignUpError("Please enter a valid full name");
      } else if (errorMessage.includes("invalid password")) {
        setSignUpError("Please enter a valid password");
      } else if (errorMessage.includes("password do not match")) {
        setSignUpError("Passwords do not match");
      } else {
        setSignUpError(errorMessage || "An unknown error occurred");
      }
    }
  };

  const handleGoogleSignUp = () => {
    window.location.assign(GOOGLE_LOGIN_URL);
  };

  return (
    <div className={styles.signUpPage}>
      <div className={styles.container}>
        <div className={styles.formSection}>
          <h1 className={styles.title}>SIGN UP</h1>
          <p className={styles.subtitle}>
            How to i get started lorem ipsum dolor at?
          </p>
          {signUpError && <div className={styles.errorText}>{signUpError}</div>}
          <form className={styles.form}>
            <div className={styles.inputGroup}>
              <span className={styles.icon}>👤</span>
              <input
                type="text"
                placeholder="Full Name"
                className={styles.input}
                value={fullName}
                onChange={(e) => setFullName(e.target.value)}
              />
            </div>

            <div className={styles.inputGroup}>
              <span className={styles.icon}>✉️</span>
              <input
                type="email"
                placeholder="Email"
                className={styles.input}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>

            <div className={styles.inputGroup}>
              <span className={styles.icon}>🔒</span>
              <input
                type="password"
                placeholder="Password"
                className={styles.input}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>

            <div className={styles.inputGroup}>
              <span className={styles.icon}>🔒</span>
              <input
                type="password"
                placeholder="Repeat Password"
                className={styles.input}
                value={repeatPassword}
                onChange={(e) => setRepeatPassword(e.target.value)}
              />
            </div>

            <Button className={styles.signUpButton} onClick={(e) => signUp(e)}>
              Sign Up Now
            </Button>
          </form>

          <p className={styles.orText}>
            or{" "}
            <Link to="/login" className={styles.link}>
              login
            </Link>
          </p>

          <p className={styles.socialHeader}>Sign up with Others</p>

          <Button className={styles.googleButton} type="button" onClick={handleGoogleSignUp}>
            <span className={styles.googleLogo}>
              <svg width="20" height="20" viewBox="0 0 24 24">
                <path
                  fill="#4285F4"
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                />
                <path
                  fill="#34A853"
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                />
                <path
                  fill="#FBBC05"
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                />
                <path
                  fill="#EA4335"
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                />
              </svg>
            </span>
            Sign up with Google
          </Button>
        </div>

        <div className={styles.imageSection}>
          <div className={styles.imageWrapper}>
            <img
              src="/sculp.jpg"
              alt="Classical sculpture"
              className={styles.image}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
