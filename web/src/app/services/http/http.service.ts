import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { take } from 'rxjs/operators';
import { environment } from 'src/environments/environment';


@Injectable({
  providedIn: 'root'
})
export class HttpService {
  protected cache = new Map();

  constructor(
    protected http: HttpClient,
  ) {
  }

  get(url: string): Promise<any> {
    return this._promise(this.http.get<any>(this._build(url), { headers: this._headers(null) }));
  }

  post(url: string, data: object): Promise<any> {
    return this._promise(this.http.post<any>(this._build(url), data, { headers: this._headers(data) }));
  }

  put(url: string, data: object): Promise<any> {
    return this._promise(this.http.put<any>(this._build(url), data, { headers: this._headers(data) }));
  }

  delete(url: string): Promise<any> {
    return this._promise(this.http.delete<any>(this._build(url), { headers: this._headers(null) }));
  }

  protected _promise(obs: Observable<any>): Promise<any> {
    return new Promise((resolve, reject) => {
      obs
        .pipe(
          take(1)
        )
        .subscribe(
          (value: any) => resolve(value),
          (err: any) => reject(err),
          () => { },
        );

    });
  }

  protected _build(url: string): string {
    return `/${environment.apiprefix}/${url}`;
  }

  protected _headers(body?: any): HttpHeaders {
    let head = new HttpHeaders();

    if (body instanceof FormData) {
      // head = head.set('Content-Type', 'multipart/form-data');
    } else {
      head = head.set('Content-Type', 'application/json');
    }

    return head;
  }

}
