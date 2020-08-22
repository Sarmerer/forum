import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { HomeComponent } from 'src/app/home/home.component';
import { ProfileComponent } from 'src/app/profile/profile.component';
import { LoginComponent } from 'src/app/login/login.component';

const routes: Routes = [
	{
		path: '',
		redirectTo: 'deals',
		pathMatch: 'full'
	},
	{
		path: 'home',
		component: HomeComponent
	},
	{
		path: 'profile',
		component: ProfileComponent
	},
	{
		path: 'login',
		component: LoginComponent
	}
];

@NgModule({
	imports: [ RouterModule.forRoot(routes) ],
	exports: [ RouterModule ]
})
export class AppRoutingModule {}
