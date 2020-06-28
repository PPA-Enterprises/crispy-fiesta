import {
  Component,
  OnInit,
  Output,
  EventEmitter,
  OnDestroy,
  ViewChild,
  ElementRef,
  Renderer2
} from "@angular/core";
import { ConfigService } from "../services/config.service";
import { Subscription } from "rxjs";
import { CustomizerService } from '../services/customizer.service';

@Component({
  selector: "app-customizer",
  templateUrl: "./customizer.component.html",
  styleUrls: ["./customizer.component.scss"]
})
export class CustomizerComponent implements OnInit, OnDestroy {
  @ViewChild("customizer") customizer: ElementRef;

  size: string;
  isBgImageDisplay: boolean = true;
  isOpen = true;
  public config: any = {};
  layoutSub: Subscription;

  constructor(
    private renderer: Renderer2,
    private configService: ConfigService,
    public customizerService: CustomizerService,
  ) {

    this.config = this.configService.templateConf;
    this.isOpen = !this.config.layout.customizer.hidden;

    if (this.config.layout.sidebar.size) {
      this.size = this.config.layout.sidebar.size;
    }
  }

  @Output() directionEvent = new EventEmitter<Object>();

  ngOnInit() {

  }


  changeSidebarWidth(value) {
    this.size = value;
    this.customizerService.changeSidebarWidth(value);
  }

  toggleCustomizer() {
    if (this.isOpen) {
      this.renderer.removeClass(this.customizer.nativeElement, "open");
      this.isOpen = false;
    } else {
      this.renderer.addClass(this.customizer.nativeElement, "open");
      this.isOpen = true;
    }
  }

  closeCustomizer() {
    this.renderer.removeClass(this.customizer.nativeElement, "open");
    this.isOpen = false;
  }

    bgImageDisplay(e) {
    if (e.target.checked) {
      this.isBgImageDisplay = true;
    } else {
      this.isBgImageDisplay = false;
    }
    //emit event to FUll Layout
    this.customizerService.bgImageDisplay(e);
  }


  ngOnDestroy() {
    if (this.layoutSub) {
      this.layoutSub.unsubscribe();
    }
  }


}
