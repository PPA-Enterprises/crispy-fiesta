import { Router } from '@angular/router';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { User } from '../models/user.model';
import { Cred } from '../models/cred.model';

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
    
    return this.http.post<User>()

    // return new Promise(function(resolve, reject) {
    //   setTimeout(function() {
    //     resolve(true);
    //   }, 1000);
    // });

  }

  logout() {
    
  }

  isAuthenticated() {
    return true;
  }
}
