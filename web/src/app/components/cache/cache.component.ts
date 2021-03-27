import { Component, OnInit } from '@angular/core';
import { Data } from '@angular/router';
import { CacheItem } from 'src/app/models/cache';
import { HttpService } from 'src/app/services/http/http.service';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-cache',
  templateUrl: './cache.component.html',
  styleUrls: ['./cache.component.scss']
})
export class CacheComponent implements OnInit {

  list: CacheItem[] = [];

  constructor(public http: HttpService) {

  }

  ngOnInit(): void {
    this.CacheList();
  }


  CacheList(): void {
    this.http.get(environment.cache_list)
      .then((value: any) => {
        this.list = value;
      })
      .catch((e: any) => {
        console.error(e);
      });
  }

}
