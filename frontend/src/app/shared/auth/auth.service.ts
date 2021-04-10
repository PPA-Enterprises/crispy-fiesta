import { Router } from '@angular/router';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { User } from '../models/user.model';
import { HttpClient } from '@angular/common/http';
import * as moment from "moment";

@Injectable()
export class AuthService {
  url ="http://ppaenterprises.com/"
  // private user: Observable<firebase.User>;
  // private userDetails: firebase.User = null;

  constructor(public router: Router, private http: HttpClient) {

  }

  signupUser(email: string, password: string) {
    //your code for signing up the new user
  }

  signinUser(email: string, password: string) {
    
    return this.http.post<any>(this.url + "api/v1/auth/", {email: email, password: password})
      .pipe(map(result => {
        console.log(result);
        localStorage.setItem('token', result.payload.token);
      }));
      
      
  }

    

  logout() {
    localStorage.removeItem('token');
  }

  public get isAuthenticated(): boolean {
    return (localStorage.getItem('token') !== null);
  }

  public tokenFromLocalStorage() : string {
    return localStorage.getItem('token');
  }
}
