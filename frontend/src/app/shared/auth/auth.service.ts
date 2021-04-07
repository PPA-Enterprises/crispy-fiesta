import { Router } from '@angular/router';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable()
export class AuthService {
  // private user: Observable<firebase.User>;
  // private userDetails: firebase.User = null;

  constructor(public router: Router) {

  }

  signupUser(email: string, password: string) {
    //your code for signing up the new user
  }

  signinUser(email: string, password: string) {
    //your code for checking credentials and getting tokens for for signing in user
    // return this._firebaseAuth.signInWithEmailAndPassword(email, password)

    //uncomment above firebase auth code and remove this temp code
    return new Promise(function(resolve, reject) {
      setTimeout(function() {
        resolve(true);
      }, 1000);
    });

  }

  logout() {

  }

  isAuthenticated() {
    return true;
  }

  public tokenFromLocalStorage() : string {
    return localStorage.getItem('token');
  }
}
