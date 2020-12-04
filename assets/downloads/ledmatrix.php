<?
/*session_start();
if(!is_array($_SESSION[data]))
{
        $_SESSION[data]=array();
        for($i=0;$i<64;$i++) $_SESSION[data][$i]=0;
}*/
if(is_array($_POST))
{
        if($_POST[cmd]=="down")
        {
                $_POST[data]= "bbbbbbbb".$_POST[data];
        }
        else if($_POST[cmd]=="up")
        {
                $ts=str_split($_POST[data],8);
                $_POST[data]="";
                for($i=1;$i<8;$i++)
                {
                        $_POST[data].=$ts[$i];
                }
                $_POST[data].="bbbbbbbb";
        }
        else if($_POST[cmd]=="right")
        {
                for($i=62;$i>-1;$i--)
                {
                        if(($i%8)!=7)$_POST[data][$i+1]=$_POST[data][$i];
                        if(($i%8)==0)$_POST[data][$i]='b';
                }
        }
        else if($_POST[cmd]=="left")
        {
                for($i=0;$i<64;$i++)
                {
                        if(($i+1)%8)$_POST[data][$i]=$_POST[data][$i+1];
                        if(($i%8)==7)$_POST[data][$i]='b';
                }
        }
        else if($_POST[cmd]=="all")
        {
                for($i=0;$i<64;$i++)$_POST[data][$i]=$_POST[color];
        }
        else if($_POST[cmd]=="none")
        {
                for($i=0;$i<64;$i++)$_POST[data][$i]='b';
        }
}

?>
<html>
<head><title>8x8 bicolor led matrix char generator v0.3</title>
<style type="text/css">
        .b {
                background: #000000;
        }
        .r {
                background: #ff0000;
        }
        .g {
                background: #00ff00;
        }
        .y {
                background: #ffff00;
        }
</style>
</head>
<body>
<script type="text/javascript">
function changeColor(nTDId)
{
        for (var i=0; i < document.form1.color.length; i++)
        {
                if (document.form1.color[i].checked)
                {
                        document.getElementById(nTDId).className = document.form1.color[i].value;
           }
        }
        updateData();
}
function updateData()
{
        document.form1.data.value = document.getElementById(0).className;
        for(var i=1; i<64; i++)
        {
                document.form1.data.value += document.getElementById(i).className;
        }
}
</script>
<form name="form1" method="post">
<? //<input type="hidden" name="data" value="bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"> ?>
<input type="hidden" name="data" value="<? echo $_POST[data] ?>">

<table><tr><td></td><td><center><input type="submit" name="cmd" value="up"></center></td><td></td></tr>
<tr><td><input type="submit" name="cmd" value="left"></td><td><table bgcolor="#aabbaa" border="1">
<?php
for($row=0;$row<8;$row++)
{
        echo "<tr>\n";
        for($col=0;$col<8;$col++)
        {
                echo '<td id="'. (($row*8)+$col) .'" class="';
                if($_POST[cmd]) echo $_POST[data][(($row*8)+$col)];
                else echo 'b';
                echo '" onClick="changeColor('.(($row*8)+$col).');">';
                echo '&nbsp&nbsp&nbsp&nbsp;</td>'."\n";
        }
        echo "</tr>\n";
}
?>

<script type="text/javascript">
updateData();
</script>

</table></td><td><input type="submit" name="cmd" value="right"></td><td>
<table>
<tr><td><input type="radio" name="color" value="b" <? if(!$_POST || $_POST[color]=='b') echo 'checked="checked"'?>>black</td></tr>
<tr><td><input type="radio" name="color" value="r" <? if($_POST[color]=='r') echo 'checked="checked"'?>>red</td></tr>
<tr><td><input type="radio" name="color" value="g" <? if($_POST[color]=='g') echo 'checked="checked"'?>>green</td></tr>
<tr><td><input type="radio" name="color" value="y" <? if($_POST[color]=='y') echo 'checked="checked"'?>>yellow</td></tr>
</table>
</td></tr>
<tr><td colspan=3><center><input type="submit" name="cmd" value="down"></center></td></tr>

<tr height=100>
<td><input type="submit" name="cmd" value="all" ></td>
<td><center><input type="submit" width=100 name="cmd" value="OK" ></center></td>
<td><input type="submit" name="cmd" value="none" ></td><td></td>
</tr></table>
</form>
<pre>
<?
$result=0;
echo '{';
for($i=0;$i<8;$i++)
{
        for($j=0;$j<8;$j++)
        {
                switch($_POST[data][($i*8)+$j])
                {
                case 'b': //echo "black";
                        $result*=4;
                        break;
                case 'g': //echo "green";
                        $result*=2; $result++; $result*=2;
                        break;
                case 'r': //echo "red"; 
                        $result*=4; $result++;
                        break;
                case 'y': //echo "yellow";
                        $result*=2; $result++; $result*=2; $result++;
                        break;
                }
//              echo "\n";
        }
        echo '0x'.dechex($result);
        if($i!=7) echo ',';
        $result=0;
}

echo '}';
echo "<br>";

//print_r($_POST);
//print_r($_GET);

?>
</pre>
</body></html>