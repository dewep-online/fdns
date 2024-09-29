import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { UXWBServicesModule } from '@uxwb/ngx-services';
import { UXWBUIModule } from '@uxwb/ngx-ui';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { PagesModule } from './pages/pages.module';

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    PagesModule,
    UXWBUIModule,
    UXWBServicesModule.forRoot({ ajaxPrefixUrl: '/api', webSocketUrl: '/ws' }),
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {
}
