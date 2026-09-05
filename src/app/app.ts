import { Component } from '@angular/core';
import { IonApp } from '@ionic/angular/ion-app';
import { IonRouterOutlet } from '@ionic/angular/ion-router-outlet';
import { addIcons } from 'ionicons';
import {
  enterOutline,
  logoApple,
  logoFacebook,
  logoGithub,
  logoGoogle,
  logoWindows,
  mailOutline,
  personAdd,
  personAddOutline,
} from 'ionicons/icons';

addIcons({
  enterOutline,
  logoApple,
  logoFacebook,
  logoGithub,
  logoGoogle,
  logoWindows,
  mailOutline,
  personAdd,
  personAddOutline,
});

@Component({
  imports: [IonApp, IonRouterOutlet],
  selector: 'app-root',
  templateUrl: './app.html',
  styleUrl: './app.scss',
})
export class App {
  protected title = 'in1app';
}
