import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';
import { CardComponent } from './card/card.component';

@NgModule({
	declarations: [ AppComponent, LoginComponent, CardComponent ],
	imports: [ BrowserModule, AppRoutingModule, FormsModule ],
	providers: [],
	bootstrap: [ AppComponent ]
})
export class AppModule {}
