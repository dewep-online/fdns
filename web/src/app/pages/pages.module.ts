import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { UXWBUIModule } from '@uxwb/ngx-ui';
import { AdblockComponent } from './adblock/adblock.component';
import { ErrorComponent } from './error/error.component';
import { HomeComponent } from './home/home.component';

@NgModule({
  declarations: [
    ErrorComponent,
    HomeComponent,
    AdblockComponent,
  ],
  imports: [
    CommonModule,
    UXWBUIModule,
  ],
  exports:[
    ErrorComponent,
    HomeComponent,
    AdblockComponent,
  ],
})
export class PagesModule { }
