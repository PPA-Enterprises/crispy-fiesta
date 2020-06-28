import { Component, OnInit, ElementRef, Inject, Renderer2, ViewChild, OnDestroy, ChangeDetectorRef, ChangeDetectionStrategy } from '@angular/core';
import { ConfigService } from 'app/shared/services/config.service';
import { DOCUMENT } from '@angular/common';
import { Subscription } from 'rxjs';
import { CustomizerService } from 'app/shared/services/customizer.service';

@Component({
  selector: 'app-content-layout',
  templateUrl: './content-layout.component.html',
  styleUrls: ['./content-layout.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})

export class ContentLayoutComponent implements OnInit, OnDestroy {
  public config: any = {};
  layoutSub: Subscription;
  @ViewChild('content-wrapper') wrapper: ElementRef;


  constructor(private configService: ConfigService,
    @Inject(DOCUMENT) private document: Document,
    private renderer: Renderer2, private cdr: ChangeDetectorRef,
    private customizerService: CustomizerService
  ) {
    this.config = this.configService.templateConf;
    this.renderer.addClass(this.document.body, "auth-page");
  }

  ngOnInit() {
    this.layoutSub = this.configService.templateConf$.subscribe((templateConf) => {
      if (templateConf) {
        this.config = templateConf;
      }
      this.loadLayout();
      this.cdr.markForCheck();

    })
  }

  loadLayout() {

    this.removeTransparentBGClasses();

    if (this.config.layout.variant === "Light") {
      this.renderer.removeClass(this.document.body, "layout-dark");
      this.renderer.removeClass(this.document.body, "layout-transparent");
    }
    else if (this.config.layout.variant === "Dark") {
      this.renderer.removeClass(this.document.body, "layout-transparent");
      this.renderer.addClass(this.document.body, "layout-dark");
    }
    else if (this.config.layout.variant === "Transparent") {
      this.renderer.addClass(this.document.body, "layout-dark");
      this.renderer.addClass(this.document.body, "layout-transparent");
      this.renderer.addClass(this.document.body, this.config.layout.sidebar.backgroundColor);
    }


    this.renderer.removeClass(this.document.body, "menu-expanded");
    this.renderer.removeClass(this.document.body, "navbar-static");
    this.renderer.removeClass(this.document.body, "menu-open");
    this.renderer.addClass(this.document.body, "blank-page");


  }

  removeTransparentBGClasses() {
    this.customizerService.transparent_colors.forEach(_ => {
      this.renderer.removeClass(this.document.body, _.class);
    });

    this.customizerService.transparent_colors_with_shade.forEach(_ => {
      this.renderer.removeClass(this.document.body, _.class);
    });
  }

  ngOnDestroy() {
    this.renderer.removeClass(this.document.body, "auth-page");
    if (this.layoutSub) {
      this.layoutSub.unsubscribe();
    }
  }

}
