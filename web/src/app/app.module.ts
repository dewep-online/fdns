import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { ServiceWorkerModule } from '@angular/service-worker';
import { UXWBServicesModule } from '@uxwb/ngx-services';
import { UXWBUIModule } from '@uxwb/ngx-ui';
import { environment } from '../environments/environment';
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
    UXWBUIModule,
    UXWBServicesModule.forRoot({ ajaxPrefixUrl: '/api', webSocketUrl: '/ws' }),
    ServiceWorkerModule.register('ngsw-worker.js', {
      enabled: environment.production,
      registrationStrategy: 'registerWhenStable:30000',
    }),
    PagesModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {
}
