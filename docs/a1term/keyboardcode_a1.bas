'----------------------------------------------------------------------------
'                PC AT-KEYBOARD for A1-REP
'                  (c) 2006 F.X.ACHATZ
'----------------------------------------------------------------------------
'connect PC AT keyboard clock to PIND.2 on the Tiny2313
'connect PC AT keyboard data to PIND.4 on the Tiny2313
'PortB0 - B6 = Data Out
'Portb,7 = STROBE (5ms Keyboard and 1ms via serial)
'Portd.3 = Reset 6502

$regfile = "Attiny2313.dat"

$crystal = 8000000
$baud = 2400

$hwstack = 32
$swstack = 16
$framesize = 32


Config Keyboard = Pind.2 , Data = Pind.4 , Keydata = Keydata
Config Portb = Output
Config Int0 = Falling                                       'generate an interrupt on the falling edge.
Config Porta.0 = Output
Set Porta.0
Config Portd.3 = Output
Set Portd.3
Config Portd.5 = Output
Set Portd.5

On Int0 Isr0
On Urxc Get_ser
Enable Int0
Enable Urxc                                                 ' RxD interrupt flag

Dim A As Byte
Dim B As Byte
Dim C As Byte
Dim Xx As Byte
Dim Ser_rx As Byte                                          ' ser rx byte
Dim Ser_received_flag As Bit                                ' flag to start routine
Dim Ser_pace As Byte

Declare Sub Toetsin
Declare Sub Controlkey_aone
Declare Sub Reset_aone
Declare Sub Clearscrn_aone
Declare Sub Ascii_out(byval Xx As Byte)


Enable Interrupts                                           'Global Interrupt Enable

Waitms 1000

B = 0
C = 0

Do
   If B > 0 Then
      Select Case B
         Case 128 : Gosub Controlkey_aone
         Case 129 : Gosub Reset_aone
         Case 130 : Gosub Clearscrn_aone
         Case Else : Gosub Toetsin
      End Select
   End If


If Ser_received_flag = 1 Then                               'a byte has been received RS232
   Disable Serial
   If Ser_rx > &H60 And Ser_rx < &H7B Then
      Ser_rx = Ser_rx And &HDF
   End If

   Select Case Ser_rx
      Case 125 : Gosub Reset_aone                           'curl close
      Case 123 : Gosub Clearscrn_aone                       'curl open
      Case Else : Call Ascii_out(ser_rx)
   End Select

   If Ser_pace > 0 Then
      Reset Portd.5
   End If

   Ser_pace = 20

   Reset Ser_received_flag                                  'reset the ser rx flag
   Enable Serial
   Ser_rx = 0
End If
Waitms 1
If Ser_pace > 0 Then
   Decr Ser_pace
   If Ser_pace = 0 Then
      Set Portd.5
   End If
End If

Loop
End

Sub Controlkey_aone
   If C = 128 Then
      C = 0
      Else
      C = B
   End If
   B = 0
End Sub


Sub Reset_aone
Reset Portd.3 : Waitms 20 : Set Portd.3 : B = 0
End Sub

Sub Clearscrn_aone
Reset Porta.0 : Waitms 20 : Set Porta.0 : B = 0
End Sub

Sub Toetsin
If C = 128 Then
   B = B And &H1F
   C = 0
End If
Call Ascii_out(b)
B = 0
End Sub


Sub Ascii_out(byval Xx As Byte)
   Portb = Xx
   sbi portb,7
   Waitms 1
   cbi portb,7
End Sub

'-------- Interrupt routine to fetch serial byte from port--------------
Get_ser:
Ser_rx = Udr
Set Ser_received_flag
Return


Isr0:                                                       'interrupt routine
B = Getatkbd()                                              'Get data from AT keyboard
Set Gifr.intf0                                              'Clear the External Interrupt Flag 0
Return


Keydata:                                                    ''This is the key translation table for the function Getatkbd().

'normal keys lower case
Data 0 , 0 , 0 , 0 , 0 , 130 , 0 , 129 , 0 , 0 , 0 , 0 , 0 , 9 , 0 , 0
Data 0 , 0 , 0 , 0 , 128 , 81 , 49 , 0 , 0 , 0 , 90 , 83 , 65 , 87 , 50 , 0
Data 0 , 67 , 88 , 68 , 69 , 52 , 51 , 0 , 0 , 32 , 86 , 70 , 84 , 82 , 53 , 0
Data 0 , 78 , 66 , 72 , 71 , 89 , 54 , 0 , 0 , 76 , 77 , 74 , 85 , 55 , 56 , 0
Data 0 , 44 , 75 , 73 , 79 , 48 , 57 , 0 , 0 , 46 , 47 , 76 , 59 , 80 , 45 , 0
Data 0 , 0 , 39 , 0 , 91 , 61 , 57 , 0 , 0 , 0 , 13 , 93 , 0 , 92 , 0 , 0
Data 0 , 62 , 0 , 0 , 0 , 8 , 223 , 0 , 49 , 49 , 52 , 52 , 55 , 0 , 0 , 0
Data 48 , 46 , 50 , 53 , 54 , 56 , 27 , 0 , 0 , 43 , 51 , 45 , 42 , 57 , 0 , 0


'shifted keys UPPER case
Data 0 , 0 , 0 , 0 , 0 , 130 , 0 , 129 , 0 , 0 , 0 , 0 , 0 , 9 , 0 , 0
Data 0 , 0 , 0 , 0 , 128 , 81 , 33 , 0 , 0 , 0 , 90 , 83 , 65 , 87 , 64 , 0
Data 0 , 67 , 88 , 68 , 69 , 36 , 35 , 0 , 0 , 0 , 86 , 70 , 84 , 82 , 37 , 0
Data 0 , 78 , 66 , 72 , 71 , 89 , 94 , 0 , 0 , 0 , 77 , 74 , 85 , 38 , 42 , 0
Data 0 , 60 , 75 , 73 , 79 , 41 , 40 , 0 , 0 , 62 , 63 , 76 , 58 , 80 , 95 , 0
Data 0 , 0 , 34 , 0 , 123 , 43 , 0 , 0 , 0 , 0 , 0 , 125 , 0 , 124 , 0 , 0
Data 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0
Data 0 , 46 , 0 , 0 , 0 , 0 , 4 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0 , 0

'Code 128 = controltoets = plaats 20
'code 129 = reset = F12 = plaats 7
'code 130 = clearscreen = F1 = plaats 5
'code 9 = tab = plaats 13