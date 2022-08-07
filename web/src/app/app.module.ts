import {ScrollingModule} from '@angular/cdk/scrolling';
import {HTTP_INTERCEPTORS} from '@angular/common/http';
import {NgModule} from '@angular/core';
import {FormsModule} from '@angular/forms';
import {BrowserModule} from '@angular/platform-browser';
import {DuiCoreModule} from '@deweppro/core';
import {ToastrModule} from 'ngx-toastr';
import {DynamicComponent} from 'src/app/components/dynamic/dynamic.component';
import {PipeModule} from 'src/app/pipe/pipe.module';
import {ErrorInterceptor} from 'src/app/services/error-interceptor.service';
import {environment} from 'src/environments/environment';
import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {AdblockComponent} from './components/adblock/adblock.component';
import {FixedComponent} from './components/fixed/fixed.component';


@NgModule({
    declarations: [
        AppComponent,
        DynamicComponent,
        AdblockComponent,
        FixedComponent
    ],
    imports: [
        BrowserModule,
        AppRoutingModule,
        FormsModule,
        PipeModule,
        DuiCoreModule.forRoot(environment.apiPrefix),
        ToastrModule.forRoot({
            preventDuplicates: true,
            progressBar: true,
            positionClass: 'toast-bottom-right'
        }),
        ScrollingModule
    ],
    providers: [
        {provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor, multi: true}
    ],
    bootstrap: [AppComponent]
})
export class AppModule {
}
