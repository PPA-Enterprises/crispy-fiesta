import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpErrorResponse
} from '@angular/common/http';
import { tap } from 'rxjs/operators';
import { AuthService } from './auth.service';
import { retry, catchError } from 'rxjs/operators';
import { Observable, throwError } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TokenInterceptorService implements HttpInterceptor {
constructor(private _authService: AuthService) {}

intercept(
    request: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    console.log(this._authService.tokenFromLocalStorage())
    request = request.clone({
      setHeaders: {
        Authorization: `Bearer ${this._authService.tokenFromLocalStorage()}`
      }
    });
    return next.handle(request).pipe(
      retry(1),
      catchError((error: HttpErrorResponse) => {
        let errorMessage = '';
        if (error.error instanceof ErrorEvent) {
          //client side error
          errorMessage = `Error: ${error.error.message}`;
        } else {
          //server side error
          errorMessage = `Error Code: ${error.status}\nMessage: ${error.message}`;
        }
        //console.log(error)

        // for displaying an error to the user
        /*if(error.status == 0) {
          this._utils.error('Internal Server Error');
        } else {
          this._utils.error(error && error.error ? error.error : '');
        }*/
        
        return throwError(errorMessage);
      })
    )
  }

}
